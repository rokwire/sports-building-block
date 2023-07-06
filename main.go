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

package main

import (
	"log"
	"os"
	"sport/driver/web"

	"github.com/rokwire/logging-library-go/v2/logs"
)

var (
	// Version : version of this executable
	Version string
)

func main() {
	if len(Version) == 0 {
		Version = "dev"
	}

	serviceID := "sports"
	loggerOpts := logs.LoggerOpts{
		SuppressRequests: logs.NewStandardHealthCheckHTTPRequestProperties(serviceID + "/version"),
	}
	logger := logs.NewLogger(serviceID, &loggerOpts)

	log.Printf("Version=%s", Version)

	//ftp credentials
	ftpHost := getEnvKey("XML_FEED_FTP_HOST")
	ftpUser := getEnvKey("XML_FEED_FTP_USER")
	ftpPassword := getEnvKey("XML_FEED_FTP_PASSWORD")

	// internal API KEY
	ssInternalAPIKey := getEnvKey("SS_INTERNAL_API_KEY")

	// host
	ssHost := getEnvKey("SS_HOST")
	coreURL := getEnvKey("SS_CORE_BB_URL")

	port := getEnvKey("PORT")
	if len(port) == 0 {
		port = "80"
	}

	///////////////////////////////////
	appID := getEnvKey("SPORTS_APP_ID")
	orgID := getEnvKey("SPORTS_ORG_ID")

	// web adapter
	webAdapter := web.NewWebAdapter(Version, port, appID, orgID, ssInternalAPIKey, ssHost, coreURL, ftpHost, ftpUser, ftpPassword, logger)
	webAdapter.Start()
	///////////////////////////////////
}

func getEnvKey(key string) string {
	//get from the environment
	value, exist := os.LookupEnv(key)
	if !exist {
		log.Fatal("No provided environment variable for " + key)
	}
	if isDevBuild() {
		log.Printf("%s=%s", key, value)
	}
	return value
}

func isDevBuild() bool {
	return Version == "dev"
}
