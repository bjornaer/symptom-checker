package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/bjornaer/sympton-checker/ent"
	"github.com/bjornaer/sympton-checker/ent/ailment"
	"github.com/bjornaer/sympton-checker/internal/symptoms"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/charset"
)

const bearerKey = "bearerKey"

type Server struct {
	client *ent.Client
	*gin.Engine
}

type AilmentResponse struct {
	Name      string `json:"name"`
	Frequency int    `json:"frequency"`
	Id        int    `json:"id"`
}

// CORSMiddleware to let us ignore cors for local development
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewServer(client *ent.Client) *Server {
	r := gin.Default()
	r.Use(CORSMiddleware())
	s := &Server{client: client, Engine: r}
	// Serve frontend static files
	r.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	api := r.Group("/api")
	{
		api.GET("/symptoms", s.GetSymptoms)
		api.POST("/symptoms", s.GetAilmentsForSymptoms)
	}
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
			}).
			AllX(c)
		filtered := filterAilments(a, hpoID)
		ailments = append(ailments, filtered...)
	}
	ailmentHistogram := constructResponse(ailments)
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

func unique(sample []*AilmentResponse) []*AilmentResponse {
	var unique []*AilmentResponse
	type key struct {
		Name          string
		Frequency, Id int
	}
	m := make(map[key]int)
	for _, v := range sample {
		k := key{v.Name, v.Frequency, v.Id}
		if i, ok := m[k]; ok {
			// Overwrite previous value per requirement in
			// question to keep last matching value.
			unique[i] = v
		} else {
			// Unique key found. Record position and collect
			// in result.
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}

func constructResponse(ailments []*ent.Ailment) []*AilmentResponse {
	ailmentList := []*AilmentResponse{}
	histogram := mapAilments(ailments)
	for _, a := range ailments {
		entry := &AilmentResponse{Name: a.Name, Id: a.ID, Frequency: histogram[a.Name]}
		ailmentList = append(ailmentList, entry)
	}
	// return the objects ordered by frequency
	sort.Slice(ailmentList, func(i, j int) bool {
		return ailmentList[i].Frequency > ailmentList[j].Frequency
	})
	return unique(ailmentList)
}
