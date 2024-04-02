package services

import (
	"errors"
	"fmt"
	"log"
)

type IJsonFormatService interface {
	ProcessJson(reqbody map[string]interface{}) (map[string]interface{}, error)
	ProcessAttributes(data map[string]interface{}) map[string]interface{}
	ProcessTraits(data map[string]interface{}) map[string]Trait
}
type JsonFormatservicestruct struct {
	Constantservice Iconstantservices
}

func JsonFormatserviceCtor(Constantservice Iconstantservices) *JsonFormatservicestruct {
	return &JsonFormatservicestruct{Constantservice: Constantservice}
}

type Trait struct {
	Value interface{} `json:"value"`
	Type  interface{} `json:"type"`
}

type AttributeStruct struct {
	Value interface{} `json:"value"`
	Type  interface{} `json:"type"`
}

// Process json
func (jsonservice *JsonFormatservicestruct) ProcessJson(reqbody map[string]interface{}) (map[string]interface{}, error) {

	log.Println("jsonformat-service starts")
	if reqbody != nil {
		ModifiedData := make(map[string]interface{})
		attributes := make(map[string]interface{})
		traits := make(map[string]interface{})
		Keyvaluepair := jsonservice.Constantservice.GetDictionary()

		for dicKey, dicVal := range Keyvaluepair {
			for key, values := range reqbody {
				if key == dicKey {
					ModifiedData[dicVal] = values
				} else if len(key) >= 3 && key[:3] == "atr" {
					attributes[key] = values

				} else if len(key) >= 4 && key[:4] == "uatr" {
					traits[key] = values
				}
			}
		}

		ProcessedAttributes := jsonservice.ProcessAttributes(attributes)
		ProcessedTraits := jsonservice.ProcessTraits(traits)
		ModifiedData["attributes"] = ProcessedAttributes
		ModifiedData["traits"] = ProcessedTraits
		log.Println("jsonformat-service ends")
		return ModifiedData, nil
	}
	return nil, errors.New("send json body")
}

// Process Attributes of given json-body
func (j *JsonFormatservicestruct) ProcessAttributes(data map[string]interface{}) map[string]interface{} {
	fmt.Println("Process Attributes starts")

	attributes := make(map[string]interface{})
	for key, value := range data {
		if len(key) > 4 && key[:4] == "atrk" {
			index := key[4:]
			atrvKey := "atrv" + index
			atrtKey := "atrt" + index
			if atrv, ok := data[atrvKey]; ok {
				if atrt, ok := data[atrtKey]; ok {
					attributes[value.(string)] = AttributeStruct{
						Value: atrv,
						Type:  atrt,
					}
				}
			}
		}
	}
	fmt.Println("Process Attributes ends")
	return attributes
}

// Process Traits of given json-body
func (j *JsonFormatservicestruct) ProcessTraits(data map[string]interface{}) map[string]Trait {
	fmt.Println("Process Traits starts")
	Traits := make(map[string]Trait)
	for dataKey, value := range data {
		if len(dataKey) > 5 && dataKey[:5] == "uatrk" {
			index := dataKey[5:]
			uatrvKey := "uatrv" + index
			uatrtKey := "uatrt" + index
			if uatrv, ok := data[uatrvKey]; ok {
				if uatrt, ok := data[uatrtKey]; ok {
					Traits[value.(string)] = Trait{
						Value: uatrv,
						Type:  uatrt,
					}
				}
			}
		}
	}
	fmt.Println("Process Traits ends")
	return Traits
}
