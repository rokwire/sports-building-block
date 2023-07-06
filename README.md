# Sports Building Block

The Sports Building Block manages sports data for the Rokwire platform.

## Documentation
The functionality provided by this application is documented in the [Wiki](https://github.com/rokwire/sports-building-block/wiki).

## Set Up

### Prerequisites

Go v1.16+

### Environment variables

The following Environment variables are supported. The service will not start unless those marked as Required are supplied.

Name|Format|Required|Description
---|---|---|---
PORT | < int > | yes | Port where this application is exposed
XML_FEED_FTP_HOST | < url > | yes | The FTP server's host for the XML feed
XML_FEED_FTP_USER | < string > | yes | The user for the FRP server of the XML feed
XML_FEED_FTP_PASSWORD | < string > | yes | The user's password for the FTP server of the XML feed
SS_INTERNAL_API_KEY | < string > | yes | The API-KEY for internal communication with the Notifications BB
SS_HOST | < url > | yes | Host for the Rokwire services
SS_CORE_BB_URL | < url > | yes | The base URL for the Core BB
SPORTS_APP_ID | < string > | yes | Application ID for Sports BB
SPORTS_ORG_ID | < string > | yes | Organization ID for Sports BB

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

### Test Application APIs

Verify the service is running as calling the get version API.

#### Call get version API

curl -X GET -i https://api-dev.rokwire.illinois.edu/sports-service/version

Response
```
2.0.0
```

## Sports API Endpoints

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

## Contributing
If you would like to contribute to this project, please be sure to read the [Contributing Guidelines](CONTRIBUTING.md), [Code of Conduct](CODE_OF_CONDUCT.md), and [Conventions](CONVENTIONS.md) before beginning.

### Secret Detection
This repository is configured with a [pre-commit](https://pre-commit.com/) hook that runs [Yelp's Detect Secrets](https://github.com/Yelp/detect-secrets). If you intend to contribute directly to this repository, you must install pre-commit on your local machine to ensure that no secrets are pushed accidentally.

```
# Install software 
$ git pull  # Pull in pre-commit configuration & baseline 
$ pip install pre-commit 
$ pre-commit install
```