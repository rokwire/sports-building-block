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

package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sport/core"
	"strconv"

	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// ApisHandler structure
type ApisHandler struct {
	app *core.Application
}

// GetVersion retrieves application version
func (a *ApisHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	version := a.app.GetVersion()
	successfulResponse(w, []byte(version))
}

// GetSports retrieves sport definitions
func (a *ApisHandler) GetSports(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	sportDefinitions, err := a.app.GetSports(l, claims.OrgID)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, "sport-definitions", nil, err, http.StatusInternalServerError, true)
	}

	data, err := json.Marshal(sportDefinitions)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, "sport-definitions", nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(data)
}

// GetNews retrieves sport news
func (a *ApisHandler) GetNews(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sports := r.URL.Query()["sport"]
	limit, el := parseLimit(r)
	if el != nil {
		log.Println(el.Error())
		http.Error(w, el.Error(), http.StatusBadRequest)
		return
	}

	news, err := a.app.GetNews(id, sports, limit)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve news. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(news) == 0 {
		successfulResponse(w, []byte("[]"))
		return
	}

	newsJSON, err := json.Marshal(news)
	if err != nil {
		errMsg := "Failed to parse news to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(newsJSON))
}

// GetCoaches retrieves coaches for a team/sport
func (a *ApisHandler) GetCoaches(w http.ResponseWriter, r *http.Request) {
	sport, err := parseSport(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coaches, err := a.app.GetCoaches(*sport)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve coaches. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(coaches) == 0 {
		successfulResponse(w, []byte("[]"))
		return
	}

	coachesJSON, err := json.Marshal(coaches)
	if err != nil {
		errMsg := "Failed to parse coaches to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(coachesJSON))
}

// GetPlayers retrieves players for a team/sport
func (a *ApisHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	sport, err := parseSport(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	players, err := a.app.GetPlayers(*sport)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve players. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(players) == 0 {
		successfulResponse(w, []byte("[]"))
		return
	}

	playersJSON, err := json.Marshal(players)
	if err != nil {
		errMsg := "Failed to parse players to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(playersJSON))
}

// GetSocialNetworks retrieves social networks
func (a *ApisHandler) GetSocialNetworks(w http.ResponseWriter, r *http.Request) {
	socNets, err := a.app.GetSocialNetworks()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve sport social networks. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(socNets) == 0 {
		successfulResponse(w, []byte("[]"))
		return
	}

	socNetsJSON, err := json.Marshal(socNets)
	if err != nil {
		errMsg := "Failed to parse social networks to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(socNetsJSON))
}

// GetGames retrieves games
func (a *ApisHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	sports := r.URL.Query()["sport"]
	id, err := parseID(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := parseDate("start", r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	endDate, err := parseDate("end", r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := parseLimit(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	games, err := a.app.GetGames(sports, id, startDate, endDate, limit)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve games. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(games) == 0 {
		successfulResponse(w, []byte("[]"))
		return
	}

	gamesJSON, err := json.Marshal(games)
	if err != nil {
		errMsg := "Failed to parse games to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(gamesJSON))
}

// GetTeamSchedule retrieves schedule for a team/sport
func (a *ApisHandler) GetTeamSchedule(w http.ResponseWriter, r *http.Request) {
	sport, err := parseSport(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schedule, err := a.app.GetTeamSchedule(*sport, year)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve team schedule. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	scheduleJSON, err := json.Marshal(schedule)
	if err != nil {
		errMsg := "Failed to parse schedule to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(scheduleJSON))
}

// GetTeamRecord retrieves schedule for a team/sport
func (a *ApisHandler) GetTeamRecord(w http.ResponseWriter, r *http.Request) {
	sport, err := parseSport(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record, err := a.app.GetTeamRecord(*sport, year)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve team record. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		errMsg := "Failed to parse record to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(recordJSON))
}

// GetLiveGames retrieves current live games
func (a *ApisHandler) GetLiveGames(w http.ResponseWriter, r *http.Request) {
	liveGames, err := a.app.GetLiveGames()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve live games. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(liveGames) == 0 {
		successfulResponse(w, []byte("[]"))
		return
	}

	var encoded []map[string]string
	for _, game := range liveGames {
		gameData := game.Encode()
		encoded = append(encoded, gameData)
	}

	result, err := json.Marshal(encoded)
	if err != nil {
		errMsg := "Failed to parse live games to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(result))
}

//GetConfig retrieves the configs
func (a *ApisHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := a.app.GetConfig()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve config. Reason: %s", err.Error())
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(config)
	if err != nil {
		errMsg := "Failed to parse config to json."
		log.Printf("%s Reason: %s", errMsg, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte(result))
}

//UpdateConfig updates the configs
func (a *ApisHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	cfgBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMsg := "failed to read request body"
		log.Printf("apis -> updateConfig: failed, reason: %s", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	err = a.app.UpdateConfig(cfgBytes)
	if err != nil {
		errMsg := "failed to update config"
		log.Printf("apis -> updateConfig: failed, reason: %s", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successfulResponse(w, []byte("Successfully updated"))
}

func parseID(r *http.Request) (*string, error) {
	ids := r.URL.Query()["id"]
	idsCount := len(ids)
	if idsCount > 1 {
		return nil, fmt.Errorf("'id' query parameter's number must be max 1 - current is [%d]", idsCount)
	}

	var id *string
	if idsCount == 1 {
		id = &ids[0]
	}
	return id, nil
}

func parseLimit(r *http.Request) (int, error) {
	limits := r.URL.Query()["limit"]
	limitsCount := len(limits)
	if limitsCount > 1 {
		return 0, fmt.Errorf("'limit' query parameter's number must be max 1 - current is [%d]", limitsCount)
	}

	var limit int
	if limitsCount == 1 {
		val, limitErr := strconv.Atoi(limits[0])
		if limitErr != nil || val < 0 {
			return 0, fmt.Errorf("'limit' parameter must be positive number - current is [%s]", limits[0])
		}
		limit = val
	}
	return limit, nil
}

func parseDate(key string, r *http.Request) (*string, error) {
	dates := r.URL.Query()[key]
	var dateVal *string
	var err error
	if dates != nil {
		if len(dates) == 1 {
			dateVal = &dates[0]
			err = validateDate(dateVal)
		} else {
			err = fmt.Errorf("please provide just one 'end' query parameter")
		}
	}
	return dateVal, err
}

func parseSport(r *http.Request) (*string, error) {
	sports := r.URL.Query()["sport"]
	if (sports == nil) || (len(sports) != 1) {
		errMsg := "please provide exactly one 'sport' query parameter"
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return &sports[0], nil
}

func parseYear(r *http.Request) (*int, error) {
	var y *int
	year := r.URL.Query()["year"]
	if len(year) == 1 {
		val, err := strconv.Atoi(year[0])
		if err != nil {
			errMsg := fmt.Sprintf("Invalid 'year' value [%s]. Please provide valid year number.", year[0])
			log.Printf("sidearm -> GetTeamSchedule: failed to parse year to int. Error: %s", err.Error())
			return nil, fmt.Errorf(errMsg)
		}
		y = &val
	} else if len(year) > 1 {
		errMsg := "please provide zero or one 'year' query parameter"
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return y, nil
}

func validateDate(date *string) error {
	if date == nil {
		return nil
	}
	// MM/dd/yyyy
	re := regexp.MustCompile("(0[1-9]|1[0-2])/(0[1-9]|1[0-9]|2[0-9]|3[0-1])/(1[89]|2[0-9])[0-9][0-9]")
	valid := re.MatchString(*date)
	if valid {
		return nil
	}
	return fmt.Errorf("provide valid date in format 'MM/dd/yyyy'")
}

func successfulResponse(w http.ResponseWriter, responseBytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

// NewApisHandler creates new instance
func NewApisHandler(app *core.Application) *ApisHandler {
	return &ApisHandler{app: app}
}
