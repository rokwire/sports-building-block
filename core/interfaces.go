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

package core

import (
	"sport/core/model"
	"sport/driven/storage"

	"github.com/rokwire/logging-library-go/v2/logs"
)

// Storage interface has to be implemented by all Storage adapters
type Storage interface {
	RegisterStorageListener(storageListener storage.Listener)
	GetSportsDefinitions(l *logs.Log, orgID string) ([]model.SportsDefinitions, error)
}

// Recipient entity
type Recipient struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// Provider interface has to be implemented by all sports providers
type Provider interface {
	GetNews(id *string, sport []string, limit int) ([]model.News, error)
	GetCoaches(sport string) ([]model.Coach, error)
	GetPlayers(sport string) ([]model.Player, error)
	GetSocialNetworks() ([]model.SportSocial, error)
	GetGames(sports []string, id *string, startDate *string, endDate *string, limit int) ([]model.Game, error)
	GetTeamSchedule(sport string, year *int) (*model.Schedule, error)
	GetTeamRecord(sport string, year *int) (*model.Record, error)
	GetLiveGames() ([]model.LiveGame, error)
	GetConfig() (map[string]interface{}, error)
	UpdateConfig(data []byte) error
}
