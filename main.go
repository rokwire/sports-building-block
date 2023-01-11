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
	"sport/driven/storage"
	"sport/driver/web"

	"github.com/rokwire/core-auth-library-go/v2/envloader"
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

	port := "80"

	///////////////////////////////////
	appID := getEnvKey("SPORTS_APP_ID")
	orgID := getEnvKey("SPORTS_ORG_ID")

	//NewStorageAdapter
	serviceID := "sports"

	loggerOpts := logs.LoggerOpts{SuppressRequests: logs.NewStandardHealthCheckHTTPRequestProperties(serviceID + "/version")}
	logger := logs.NewLogger(serviceID, &loggerOpts)
	envLoader := envloader.NewEnvLoader(Version, logger)

	// mongoDB adapter
	mongoDBAuth := envLoader.GetAndLogEnvVar("SPORTS_MONGO_AUTH", true, true)
	mongoDBName := envLoader.GetAndLogEnvVar("SPORTS_MONGO_DATABASE", true, false)
	mongoTimeout := envLoader.GetAndLogEnvVar("SPORTS_MONGO_TIMEOUT", false, false)

	storageAdapter := storage.NewStorageAdapter(mongoDBAuth, mongoDBName, mongoTimeout, logger)
	err := storageAdapter.Start()
	if err != nil {
		logger.Fatalf("Cannot start the mongoDB adapter: %v", err)
	}
	//NewApplication
	// web adapter
	webAdapter := web.NewWebAdapter(Version, port, appID, orgID, ssInternalAPIKey, ssHost, coreURL, ftpHost, ftpUser, ftpPassword, storageAdapter)
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
