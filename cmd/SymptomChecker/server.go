package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/bjornaer/sympton-checker/ent"
	"github.com/bjornaer/sympton-checker/ent/ailment"
	"github.com/bjornaer/sympton-checker/internal/symptoms"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/charset"
)

const bearerKey = "bearerKey"

type Server struct {
	client *ent.Client
	*gin.Engine
}

func NewServer(client *ent.Client) *Server {
	r := gin.Default()
	s := &Server{client: client, Engine: r}
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This is the symptom checker app")
	})
	r.GET("/symptoms", s.GetSymptoms)
	r.POST("/symptoms", s.GetAilmentsForSymptoms)
	return s
}

func (s *Server) LoadSymptomsFromRemote(source string) {
	resp, err := http.Get(source)
	if err != nil {
		log.Panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	var xmlContent symptoms.JDBOR // parse XML file data into here
	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&xmlContent)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("found %s ailments...\tPopulating DB...", xmlContent.HPODisorderSetStatusList.Count)
	symptoms.PopulateStore(context.Background(), s.client, &xmlContent.HPODisorderSetStatusList)
}

func (s *Server) GetAilmentsForSymptoms(c *gin.Context) {
	var symptomHPOIdList []string
	if err := c.BindJSON(&symptomHPOIdList); err != nil {
		c.String(400, "Incorrect payload")
	}
	ailments := []*ent.Ailment{}
	for _, hpoID := range symptomHPOIdList {
		a := s.client.Ailment.Query().
			Where(func(s *sql.Selector) {
				s.
					Where(sqljson.ValueContains(ailment.FieldHpos, hpoID))
				//  Where(sqljson.ValueEQ(ailment.FieldSymptoms, symptoms.FrequencyHigh, sqljson.Path("frequency")))
			}).
			AllX(c)
		filtered := filterAilments(a, hpoID)
		ailments = append(ailments, filtered...)
	}
	ailmentHistogram := mapAilments(ailments)
	c.JSON(200, ailmentHistogram)
}

func (s *Server) GetSymptoms(c *gin.Context) {
	symp := s.client.Symptom.Query().AllX(c)
	c.JSON(200, symp)
}

func filterAilments(ailments []*ent.Ailment, hpo string) []*ent.Ailment {
	filtered := []*ent.Ailment{}
	for _, a := range ailments {
		if a.Symptoms[hpo].Frequency == symptoms.FrequencyHigh {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func mapAilments(ailments []*ent.Ailment) map[string]int {
	histogram := map[string]int{}
	for _, a := range ailments {
		if val, ok := histogram[a.Name]; ok {
			histogram[a.Name] = val + 1
		} else {
			histogram[a.Name] = 1
		}
	}
	return histogram
}
