package symptoms

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"

	"github.com/bjornaer/sympton-checker/ent"
	"github.com/bjornaer/sympton-checker/ent/schema"
)

const (
	FrequencyHigh = "Very frequent (99-80%)"
	FrequencyMid  = "Frequent (79-30%)"
	FrequencyLow  = "Occasional (29-5%)"
)

type JDBOR struct {
	XMLName      xml.Name `xml:"JDBOR"`
	Text         string   `xml:",chardata"`
	Date         string   `xml:"date,attr"`
	Version      string   `xml:"version,attr"`
	Copyright    string   `xml:"copyright,attr"`
	Dbserver     string   `xml:"dbserver,attr"`
	Availability struct {
		Text    string `xml:",chardata"`
		Licence struct {
			Text            string `xml:",chardata"`
			FullName        Name   `xml:"FullName"`
			ShortIdentifier string `xml:"ShortIdentifier"`
			LegalCode       string `xml:"LegalCode"`
		} `xml:"Licence"`
	} `xml:"Availability"`
	HPODisorderSetStatusList HPODisorderSetStatusList
}

type HPODisorderSetStatusList struct {
	XMLName              xml.Name               `xml:"HPODisorderSetStatusList"`
	Text                 string                 `xml:",chardata"`
	Count                string                 `xml:"count,attr"`
	HPODisorderSetStatus []HPODisorderSetStatus `xml:"HPODisorderSetStatus"`
}

type HPODisorderSetStatus struct {
	Text             string   `xml:",chardata"`
	ID               string   `xml:"id,attr"`
	Disorder         Disorder `xml:"Disorder"`
	Source           string   `xml:"Source"`
	ValidationStatus string   `xml:"ValidationStatus"`
	Online           string   `xml:"Online"`
	ValidationDate   string   `xml:"ValidationDate"`
}

type Disorder struct {
	Text                       string                     `xml:",chardata"`
	ID                         string                     `xml:"id,attr"`
	OrphaCode                  string                     `xml:"OrphaCode"`
	ExpertLink                 ExpertLink                 `xml:"ExpertLink"`
	Name                       Name                       `xml:"Name"`
	DisorderType               DisorderType               `xml:"DisorderType"`
	DisorderGroup              DisorderGroup              `xml:"DisorderGroup"`
	HPODisorderAssociationList HPODisorderAssociationList `xml:"HPODisorderAssociationList"`
}

type ExpertLink struct {
	Text string `xml:",chardata"`
	Lang string `xml:"lang,attr"`
}

type Name struct {
	Text string `xml:",chardata"`
	Lang string `xml:"lang,attr"`
}

type DisorderType struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name Name   `xml:"Name"`
}

type DisorderGroup struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name Name   `xml:"Name"`
}

type HPODisorderAssociationList struct {
	Text                   string `xml:",chardata"`
	Count                  string `xml:"count,attr"`
	HPODisorderAssociation []HPODisorderAssociation
}

type HPODisorderAssociation struct {
	Text               string       `xml:",chardata"`
	ID                 string       `xml:"id,attr"`
	HPO                HPO          `xml:"HPO"`
	HPOFrequency       HPOFrequency `xml:"HPOFrequency"`
	DiagnosticCriteria string       `xml:"DiagnosticCriteria"`
}

type HPO struct {
	Text    string `xml:",chardata"`
	ID      string `xml:"id,attr"`
	HPOId   string `xml:"HPOId"`
	HPOTerm string `xml:"HPOTerm"`
}

type HPOFrequency struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name Name   `xml:"Name"`
}

func CreateSymptom(ctx context.Context, client *ent.Client, symptomXML HPODisorderAssociation) *schema.SymptomDetails {
	name := symptomXML.HPO.HPOTerm
	hpoID := symptomXML.HPO.HPOId
	freq := symptomXML.HPOFrequency.Name.Text
	client.Symptom.Create().SetHpo(hpoID).SetName(name).Save(ctx)
	return &schema.SymptomDetails{Name: name, HPOId: hpoID, Frequency: freq}
}

func CreateAssociatedSymptomList(ctx context.Context, client *ent.Client, symptomsList HPODisorderAssociationList) (map[string]schema.SymptomDetails, []string) {
	symptoms := map[string]schema.SymptomDetails{}
	hpos := []string{}
	for _, s := range symptomsList.HPODisorderAssociation {
		symptom := CreateSymptom(ctx, client, s)
		symptoms[symptom.HPOId] = *symptom
		hpos = append(hpos, symptom.HPOId)
	}
	return symptoms, hpos
}

func CreateAilment(ctx context.Context, client *ent.Client, ailmentXML *HPODisorderSetStatus) (*ent.Ailment, error) {
	id, err := strconv.Atoi(ailmentXML.Disorder.OrphaCode)
	if err != nil {
		return nil, fmt.Errorf("failed creating ailment: %w", err)
	}
	name := ailmentXML.Disorder.Name.Text
	associatedSymptoms, hpos := CreateAssociatedSymptomList(ctx, client, ailmentXML.Disorder.HPODisorderAssociationList)
	a, err := client.Ailment.
		Create().
		SetID(id).
		SetName(name).
		SetSymptoms(associatedSymptoms).
		SetHpos(hpos).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	// log.Println("ailment was created: ", a.Name)
	return a, nil
}

func PopulateStore(ctx context.Context, client *ent.Client, ailmentListXML *HPODisorderSetStatusList) {
	for _, a := range ailmentListXML.HPODisorderSetStatus {
		CreateAilment(ctx, client, &a)
	}
	log.Print("DB populated")
}
