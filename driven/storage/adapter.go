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
	"log"
	"sport/core/model"
	"strconv"
	"time"

	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Adapter implements the Storage interface
type Adapter struct {
	db     *database
	logger *logs.Logger
}

// Start starts the storage
func (sa *Adapter) Start() error {
	//start db
	err := sa.db.start()
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionInitialize, "storage adapter", nil, err)
	}

	//register storage listener
	sl := storageListener{adapter: sa}
	sa.RegisterStorageListener(&sl)

	return err
}

// RegisterStorageListener registers a data change listener with the storage adapter
func (sa *Adapter) RegisterStorageListener(storageListener Listener) {
	sa.db.listeners = append(sa.db.listeners, storageListener)
}

// GetSportsDefinitions retrieves sport definitions
func (sa *Adapter) GetSportsDefinitions(l *logs.Log, orgID string) ([]model.SportsDefinitions, error) {
	filter := bson.D{}
	var sports []model.SportsDefinitions
	err := sa.db.sportDefinitions.Find(filter, &sports, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, "sport-definitions", nil, err)
	}
	return sports, nil
}

// NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string, logger *logs.Logger) *Adapter {
	timeoutInt, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		logger.Warn("Setting default Mongo timeout - 1000")
		timeoutInt = 1000
	}
	timeout := time.Millisecond * time.Duration(timeoutInt)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeout, logger: logger}
	return &Adapter{db: db, logger: logger}
}

type storageListener struct {
	adapter *Adapter
}

// Listener represents storage listener
type Listener interface {
}

func abortTransaction(sessionContext mongo.SessionContext) {
	err := sessionContext.AbortTransaction(sessionContext)
	if err != nil {
		log.Printf("error on aborting a transaction - %s", err)
	}
}
