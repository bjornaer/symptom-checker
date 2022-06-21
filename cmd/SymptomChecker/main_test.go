package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	p "github.com/bjornaer/sympton-checker/cmd/SymptomChecker"
	"github.com/bjornaer/sympton-checker/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	router *p.Server
}

func (s *UnitTestSuite) SetupTest() {
	// s.dbClient, _ = db.InitDBClient()
}

func (s *UnitTestSuite) BeforeTest(suiteName, testName string) {
	// dbClient := s.dbClient
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	s.router = p.NewServer(client)
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
}

func (s *UnitTestSuite) Test_TableTest() {

	type testCase struct {
		name             string
		endpoint         string
		requestMethod    string
		expectedResponse interface{}
		payload          string
	}

	testCases := []testCase{
		{
			name:             "Hello Endpoint",
			endpoint:         "/",
			requestMethod:    http.MethodGet,
			payload:          "",
			expectedResponse: "Hello Sailor! Welcome to the Port Domain Service!",
		},
	}

	for _, testCase := range testCases {

		s.Run(testCase.name, func() {
			router := s.router
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(testCase.requestMethod, testCase.endpoint, strings.NewReader(testCase.payload))
			router.ServeHTTP(w, req)

			assert.Equal(s.T(), http.StatusOK, w.Code)
			assert.Equal(s.T(), testCase.expectedResponse, w.Body.String())

		})
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
