# Sports Service

Go project to provide rest service for rokwire sports building block.

## Set Up

### Prerequisites

Go v1.16+

### Environment variables

The following Environment variables are supported. The service will not start unless those marked as Required are supplied.

Name|Value|Required|Description
---|---|---|---
XML_FEED_FTP_HOST | < value > | yes | The ftp server's host of the xml feed
XML_FEED_FTP_USER | < value > | yes | The user for the ftp server of the xml feed
XML_FEED_FTP_PASSWORD | < value > | yes | The user's password for the ftp server of the xml feed
SS_INTERNAL_API_KEY | < value > | yes | The API-KEY for internal communication with the Notifications BB
SS_HOST | < value > | yes | The host which the service is deployed on
SS_CORE_BB_URL | < value > | yes | The base URL for the Core BB

#### Run locally without Docker

1. Clone this repo (outside GOPATH)

2. Open the terminal and go to the root folder

```bash
cd sports-service
```

3. Make the project

```bash
$ make
▶ go mod vendor
▶ Log info…
...
▶ Checking formatting…
▶ running golint…
▶ building executable(s)… 1.0.247 2021-10-12T08:56:56+0300
sport
```

4. Run the executable

```bash
./bin/sport
```

#### Run locally as Docker container

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Create Docker image

```bash
docker build -t sports-service .
```

4. Run as Docker container

```bash
docker run -e ROKWIRE_API_KEYS -e XML_FEED_FTP_HOST -e XML_FEED_FTP_USER -e XML_FEED_FTP_PASSWORD -e FIREBASE_PROJECT_ID -e FIREBASE_AUTH -e SS_FIREBASE_PROJECT_ID_SAFER -e SS_FIREBASE_AUTH_SAFER -e SS_INTERNAL_API_KEY -e SS_HOST -d -p 80:80  sports-service

docker stop $(docker ps -a -q)
```

#### Tools

##### Run golint

```bash
make lint
```

##### Run gofmt to check formatting on all source files

```bash
make checkfmt
```

##### Run gofmt to fix formatting on all source files

```bash
make fixfmt
```

##### Cleanup everything

```bash
make clean
```

## Sports API End Points

Name|Deprecated|Description
---|---|---
/sports-service/version | no | get server version
/sports-service/api/v2/config | no | get/update live games config
/sports-service/api/v2/sports | no | get sport definitions
/sports-service/api/v2/news | no | get news
/sports-service/api/v2/coaches | no | get coaches
/sports-service/api/v2/players | no | get players
/sports-service/api/v2/social | no | get social media accounts
/sports-service/api/v2/games | no | get games
/sports-service/api/v2/team-schedule | no | get team schedule
/sports-service/api/v2/team-record | no | get team record
/sports-service/api/v2/live-games | no | get current live games
