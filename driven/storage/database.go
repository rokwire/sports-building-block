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
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"sport/core/model"
	"time"

	"github.com/rokwire/logging-library-go/v2/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	logger *logs.Logger

	db       *mongo.Database
	dbClient *mongo.Client

	sportDefinitions *collectionWrapper

	listeners []Listener
}

func (m *database) start() error {
	m.logger.Info("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)

	sportDefinitions := &collectionWrapper{database: m, coll: db.Collection("sport-definitions")}

	err = m.applySportDefinitionsChecks(sportDefinitions)
	if err != nil {
		return err
	}

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	m.sportDefinitions = sportDefinitions

	//apply sport definitions
	err = m.setSportDefinitionsData(client, sportDefinitions)
	if err != nil {
		return err
	}

	return nil
}

// set sport definitions data
func (m *database) setSportDefinitionsData(client *mongo.Client,
	sportDefinitions *collectionWrapper) error {

	//get sportDefinition data
	_, err := m.getSportDefinitions(client, sportDefinitions, 19)
	if err != nil {
		m.logger.Errorf("error migraiting sport definitions data - %s", err)
		return err
	}

	m.logger.Info("setSportDefinitionsData finished")
	return nil
}

func (m *database) getSportDefinitions(client *mongo.Client, sportDefinition *collectionWrapper, count int) ([]model.SportsDefinitions, error) {
	fileBytes, err := ioutil.ReadFile("driven/storage/sport-definitions.json")
	if err != nil {
		log.Printf("Failed to read sport-definitions.json file. Reason: %s", err.Error())
		return nil, nil // the "zero" value for strings is empty string
	}

	var sdef []model.SportsDefinitions
	err = json.Unmarshal([]byte(fileBytes), &sdef)

	err = m.sportDefinitions.Find(nil, &sdef, nil)
	if err != nil {
		return nil, nil
	}
	sDCount := len(sdef)
	if sDCount == 0 {
		var insertSportDef []model.SportsDefinitions
		err = json.Unmarshal([]byte(fileBytes), &insertSportDef)
		sport := make([]interface{}, len(insertSportDef))
		for i, sd := range insertSportDef {
			sport[i] = sd
		}
		_, err = sportDefinition.InsertMany(sport, nil)
	} else {
		log.Printf("Sport definitions data is already migrated")
		return nil, nil
	}

	return nil, nil
}

func (m *database) applySportDefinitionsChecks(sportDefinitions *collectionWrapper) error {
	m.logger.Info("apply sport definitions checks.....")

	//add org id index
	err := sportDefinitions.AddIndex(bson.D{primitive.E{Key: "org_id", Value: 1}}, false, false)
	if err != nil {
		return err
	}

	m.logger.Info("accounts sport definitions passed")
	return nil
}
