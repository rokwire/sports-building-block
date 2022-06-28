package storage

import (
	"io/ioutil"
	"log"
)

// Adapter implements Storage interface
type Adapter struct {
}

// GetSports retrieves sport definitions
func (sa *Adapter) GetSports() string {
	fileBytes, err := ioutil.ReadFile("driven/storage/sport-definitions.json")
	if err != nil {
		log.Printf("Failed to read sport-definitions.json file. Reason: %s", err.Error())
		return "" // the "zero" value for strings is empty string
	}
	return string(fileBytes)
}

// NewStorageAdapter creates new instance
func NewStorageAdapter() *Adapter {
	return &Adapter{}
}
