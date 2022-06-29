// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
