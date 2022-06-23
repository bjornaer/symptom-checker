package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	p "github.com/bjornaer/sympton-checker/cmd/SymptomChecker"
	"github.com/bjornaer/sympton-checker/ent"
	"github.com/bjornaer/sympton-checker/ent/enttest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	router *p.Server
	client *ent.Client
}

func (s *UnitTestSuite) SetupTest() {
	// s.dbClient, _ = db.InitDBClient()
	client := enttest.Open(s.T(), "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	s.client = client
	s.router = p.NewServer(client)
	s.router.LoadSymptomsFromRemote("http://www.orphadata.org/data/xml/en_product4.xml")
}

func (s *UnitTestSuite) BeforeTest(suiteName, testName string) {
	// client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	// s.client = client
	// client.Schema.Create(context.Background())
	// s.router = p.NewServer(client)
	// s.router.LoadSymptomsFromRemote("http://www.orphadata.org/data/xml/en_product4.xml")
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	// s.client.Close()
}

func (s *UnitTestSuite) Test_TableTest() {

	type testCase struct {
		name             string
		endpoint         string
		requestMethod    string
		expectedStatus   int
		expectedResponse interface{}
		payload          string
	}

	testCases := []testCase{
		{
			name:             "Entry Endpoint",
			endpoint:         "/",
			requestMethod:    http.MethodGet,
			payload:          "",
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "",
		},
		{
			name:             "Get Symptoms Endpoint",
			endpoint:         "/api/symptoms",
			requestMethod:    http.MethodGet,
			expectedStatus:   http.StatusOK,
			payload:          "",
			expectedResponse: "HP:0000256",
		},
		{
			name:             "Get Ailment from Symptoms Endpoint",
			endpoint:         "/api/symptoms",
			requestMethod:    http.MethodPost,
			expectedStatus:   http.StatusOK,
			payload:          "[\"HP:0000256\"]",
			expectedResponse: "Acrocallosal syndrome",
		},
	}

	for _, testCase := range testCases {

		s.Run(testCase.name, func() {
			router := s.router
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(testCase.requestMethod, testCase.endpoint, strings.NewReader(testCase.payload))
			router.ServeHTTP(w, req)

			assert.Equal(s.T(), testCase.expectedStatus, w.Code)
			assert.Contains(s.T(), w.Body.String(), testCase.expectedResponse)

		})
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
