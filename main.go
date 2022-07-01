package main

import (
	"log"
	"os"
	"sport/driver/web"
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
	// web adapter
	webAdapter := web.NewWebAdapter(Version, port, ssInternalAPIKey, ssHost, coreURL, ftpHost, ftpUser, ftpPassword)
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
