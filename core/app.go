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
	"sport/driven/provider/sidearm"
	"sport/driven/storage"
)

// Application structure
type Application struct {
	version  string
	storage  storage.Adapter
	provider Provider
}

// GetVersion returns application's version
func (app *Application) GetVersion() string {
	return app.version
}

// GetSports retrieves sport definitions
func (app *Application) GetSports(orgID string, appID string) string {
	return app.storage.GetSports(orgID, appID)
}

// GetNews retrieves sport news
func (app *Application) GetNews(id *string, sports []string, limit int, orgID string, appID string) ([]model.News, error) {
	return app.provider.GetNews(id, sports, limit, orgID, appID)
}

// GetCoaches retrieves the coaches for specific sport
func (app *Application) GetCoaches(sport string, orgID string, appID string) ([]model.Coach, error) {
	return app.provider.GetCoaches(sport, orgID, appID)
}

// GetPlayers retrieves the players for specific sport
func (app *Application) GetPlayers(sport string, orgID string, appID string) ([]model.Player, error) {
	return app.provider.GetPlayers(sport, orgID, appID)
}

// GetSocialNetworks retrieves the social accounts
func (app *Application) GetSocialNetworks(orgID string, appID string) ([]model.SportSocial, error) {
	return app.provider.GetSocialNetworks(orgID, appID)
}

// GetGames retrieves games based on selected filters
func (app *Application) GetGames(sports []string, id *string, startDate *string, endDate *string, limit int, orgID string, appID string) ([]model.Game, error) {
	return app.provider.GetGames(sports, id, startDate, endDate, limit, orgID, appID)
}

// GetTeamSchedule retrieves the schedule for sport in a specific year
func (app *Application) GetTeamSchedule(sport string, year *int, orgID string, appID string) (*model.Schedule, error) {
	return app.provider.GetTeamSchedule(sport, year, orgID, appID)
}

// GetTeamRecord retrieves the record for a sport team
func (app *Application) GetTeamRecord(sport string, year *int, orgID string, appID string) (*model.Record, error) {
	return app.provider.GetTeamRecord(sport, year, orgID, appID)
}

// GetLiveGames retrieves details for current live games
func (app *Application) GetLiveGames(orgID string, appID string) ([]model.LiveGame, error) {
	return app.provider.GetLiveGames(orgID, appID)
}

// GetConfig retrieves provider's config
func (app *Application) GetConfig() (map[string]interface{}, error) {
	return app.provider.GetConfig()
}

// UpdateConfig updates provider's config
func (app *Application) UpdateConfig(cfgBytes []byte) error {
	return app.provider.UpdateConfig(cfgBytes)
}

// NewApplication creates new Application instance
func NewApplication(version string, internalAPIKey string, appID string, orgID string, host string, ftpHost string, ftpUser string, ftpPassword string) *Application {
	sa := storage.NewStorageAdapter()

	// Here we define current sport provider!
	sp := sidearm.NewProvider(internalAPIKey, host, ftpHost, ftpUser, ftpPassword, appID, orgID)
	sp.Start()

	return &Application{version: version, storage: *sa, provider: sp}
}
