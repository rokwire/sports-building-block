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
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rokwire/logging-library-go/v2/logs"
	"log"
	"net/http"
	"sport/core"
)

// Adapter structure
type Adapter struct {
	port string
	apis *ApisHandler
	auth *auth

	logger *logs.Logger
}

// Start adapter
func (we Adapter) Start() {

	router := mux.NewRouter().StrictSlash(true)
	defaultSubRouter := router.PathPrefix("/sports-service").Subrouter()
	apiSubRouter := defaultSubRouter.PathPrefix("/api").Subrouter()

	//////////////////////////////////////////////////
	/// General Usage APIs
	defaultSubRouter.HandleFunc("/version", we.apis.GetVersion).Methods("GET")

	//////////////////////////////////////////////////
	/// V2 APIs
	v2SubRouter := apiSubRouter.PathPrefix("/v2").Subrouter()
	v2SubRouter.HandleFunc("/config", we.corePermissionWrapFunc(we.apis.GetConfig)).Methods("GET")
	v2SubRouter.HandleFunc("/config", we.corePermissionWrapFunc(we.apis.UpdateConfig)).Methods("PUT")
	v2SubRouter.HandleFunc("/sports", we.coreWrapFunc(we.apis.GetSports)).Methods("GET")
	v2SubRouter.HandleFunc("/news", we.coreWrapFunc(we.apis.GetNews)).Methods("GET")
	v2SubRouter.HandleFunc("/coaches", we.coreWrapFunc(we.apis.GetCoaches)).Methods("GET")
	v2SubRouter.HandleFunc("/players", we.coreWrapFunc(we.apis.GetPlayers)).Methods("GET")
	v2SubRouter.HandleFunc("/social", we.coreWrapFunc(we.apis.GetSocialNetworks)).Methods("GET")
	v2SubRouter.HandleFunc("/games", we.coreWrapFunc(we.apis.GetGames)).Methods("GET")
	v2SubRouter.HandleFunc("/team-schedule", we.coreWrapFunc(we.apis.GetTeamSchedule)).Methods("GET")
	v2SubRouter.HandleFunc("/team-record", we.coreWrapFunc(we.apis.GetTeamRecord)).Methods("GET")
	v2SubRouter.HandleFunc("/live-games", we.coreWrapFunc(we.apis.GetLiveGames)).Methods("GET")
	//////////////////////////////////////////////////

	err := http.ListenAndServe(":"+we.port, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (we Adapter) coreWrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logObj := we.logger.NewRequestLog(r)
		logObj.RequestReceived()

		err := we.auth.coreAuthCheck(w, r)

		if err != nil {
			errMsg := fmt.Sprintf("Unauthorized: %s", err.Error())
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		handler(w, r)
		logObj.RequestComplete()
	}
}

func (we Adapter) corePermissionWrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logObj := we.logger.NewRequestLog(r)
		logObj.RequestReceived()

		err := we.auth.corePermissionAuthCheck(w, r)

		if err != nil {
			errMsg := fmt.Sprintf("Unauthorized: %s", err.Error())
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		handler(w, r)
		logObj.RequestComplete()
	}
}

func logRequest(req *http.Request) {
	if req == nil {
		return
	}

	method := req.Method
	path := req.URL.Path

	val, ok := req.Header["User-Agent"]
	if ok && len(val) != 0 && val[0] == "ELB-HealthChecker/2.0" {
		return
	}

	header := make(map[string][]string)
	for key, value := range req.Header {
		var logValue []string
		// Do not log api key, cookies and Authorization headers
		if (key == "Rokwire-Api-Key") || (key == "Cookie") || (key == "Authorization") {
			logValue = append(logValue, "---")
		} else {
			logValue = value
		}
		header[key] = logValue
	}
	log.Printf("%s %s %s", method, path, header)
}

// NewWebAdapter creates new instance
func NewWebAdapter(version string, port string, appID string, orgID string, internalAPIKey string, host string, coreURL string, ftpHost string, ftpUser string, ftpPassword string, logger *logs.Logger) Adapter {
	app := core.NewApplication(version, internalAPIKey, appID, orgID, host, ftpHost, ftpUser, ftpPassword)
	apis := NewApisHandler(app)
	auth := newAuth(host, coreURL)
	return Adapter{port: port, apis: apis, auth: auth, logger: logger}
}
