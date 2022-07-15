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

package source

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"sport/core/model"
	sidearmModel "sport/driven/provider/sidearm/model"
	"strconv"
	"strings"
	"time"
)

type xmlFeedFootballGame struct {
	gameID        int
	sport         string
	hasStarted    bool
	isComplete    bool
	homeScore     int
	visitingScore int
	customData    string

	//we need to support old version as well.
	clockSeconds int
	period       int
}

func (game *xmlFeedFootballGame) GetType() string {
	return game.sport
}

func (game *xmlFeedFootballGame) GetGameID() int {
	return game.gameID
}

func (game *xmlFeedFootballGame) GetPath() string {
	return game.sport
}

func (game *xmlFeedFootballGame) GetHasStarted() bool {
	return game.hasStarted
}

func (game *xmlFeedFootballGame) GetIsComplete() bool {
	return game.isComplete
}

func (game *xmlFeedFootballGame) GetClockSeconds() int {
	//we need to support old version as well.
	return game.clockSeconds
}

func (game *xmlFeedFootballGame) GetPeriod() int {
	//we need to support old version as well.
	return game.period
}

func (game *xmlFeedFootballGame) GetHomeScore() int {
	return game.homeScore
}

func (game *xmlFeedFootballGame) GetVisitingScore() int {
	return game.visitingScore
}

func (game *xmlFeedFootballGame) GetCustomData() string {
	return game.customData
}

func (game *xmlFeedFootballGame) Encode() map[string]string {

	data := make((map[string]string))
	data["Type"] = game.GetType()
	data["GameId"] = strconv.Itoa(game.GetGameID())
	data["Path"] = game.GetPath()
	data["HasStarted"] = strconv.FormatBool(game.GetHasStarted())
	data["IsComplete"] = strconv.FormatBool(game.GetIsComplete())
	data["ClockSeconds"] = strconv.Itoa(game.GetClockSeconds())
	data["Period"] = strconv.Itoa(game.GetPeriod())
	data["HomeScore"] = strconv.Itoa(game.GetHomeScore())
	data["VisitingScore"] = strconv.Itoa(game.GetVisitingScore())
	data["Custom"] = game.GetCustomData()

	return data
}

type footballCustomData struct {
	Possession string
	LastPlay   string
	Clock      string
	Phase      string
}

type xmlFootballGame struct {
	XMLName   xml.Name          `xml:"fbgame"`
	Generated string            `xml:"generated,attr"`
	Venue     xmlFootballVenue  `xml:"venue"`
	Plays     xmlFootballPlays  `xml:"plays"`
	Scores    xmlFootballScores `xml:"scores"`
}

type xmlFootballVenue struct {
	XMLName xml.Name `xml:"venue"`
	Date    string   `xml:"date,attr"`
}

type xmlFootballScores struct {
	XMLName    xml.Name           `xml:"scores"`
	ScoresList []xmlFootballScore `xml:"score"`
}

type xmlFootballScore struct {
	XMLName xml.Name `xml:"score"`
	Quarter string   `xml:"qtr,attr"`
	Vscore  string   `xml:"vscore,attr"`
	Hscore  string   `xml:"hscore,attr"`
}

type xmlFootballPlays struct {
	XMLName  xml.Name             `xml:"plays"`
	Quarters []xmlFootballQuarter `xml:"qtr"`
	Downtogo xmlFootballDowntogo  `xml:"downtogo"`
}

type xmlFootballQuarter struct {
	XMLName xml.Name              `xml:"qtr"`
	Number  string                `xml:"number,attr"`
	Plays   []xmlFootballPlay     `xml:"play"`
	Scores  []xmlFootballQtrScore `xml:"score"`
}

type xmlFootballPlay struct {
	XMLName xml.Name `xml:"play"`
	Score   string   `xml:"score,attr"`
	Vscore  string   `xml:"vscore,attr"`
	Hscore  string   `xml:"hscore,attr"`
	Clock   string   `xml:"clock,attr"`
}

type xmlFootballQtrScore struct {
	XMLName xml.Name `xml:"score"`
	Final   string   `xml:"final,attr"`
	V       string   `xml:"V,attr"`
	H       string   `xml:"H,attr"`
}

type xmlFootballDowntogo struct {
	XMLName  xml.Name `xml:"downtogo"`
	Hasball  string   `xml:"hasball,attr"`
	Qtr      string   `xml:"qtr,attr"`
	Clock    string   `xml:"clock,attr"`
	LastPlay string   `xml:"lastplay,attr"`
}

type xmlFootballSource struct {
	config  Config
	ftpConn ftpConn
}

func newXMLFootballSource(config Config, ftpHost string, ftpUser string, ftpPassword string) xmlFootballSource {
	var xmlFootballSource xmlFootballSource
	xmlFootballSource.config = config
	xmlFootballSource.ftpConn = newFTPConn(ftpHost, ftpUser, ftpPassword)
	return xmlFootballSource
}

func (xmlFootballSource *xmlFootballSource) updateConfig(config Config) {
	log.Println("xmlfootball: UpdateConfig -> config updated in xml footbal source")
	xmlFootballSource.config = config
}

func (xmlFootballSource *xmlFootballSource) load(item *sidearmModel.LiveGameItem) (model.LiveGame, error) {
	//1. load the xml data
	xmlData, err := xmlFootballSource.ftpConn.load("/Rokwire_FTP/football")
	if err != nil {
		return nil, err
	}

	//2. unmarshal
	var xmlFootballGame *xmlFootballGame
	err = xml.Unmarshal(xmlData, &xmlFootballGame)
	if err != nil {
		return nil, err
	}

	//printXMLFootballGame(xmlFootballGame)

	//3. check if the dowloaded xml is for this game. The only way we could check is to compare the dates!
	if xmlFootballSource.config.GetFootballDateCheck() && !xmlFootballSource.isForGame(xmlFootballGame, item) {
		return nil, errors.New("xmlfootball: loadFromXML -> the xml is not for this game")
	}

	//4. construct xmlFeedFootballGame
	var xmlFeedGame xmlFeedFootballGame
	xmlFeedGame.gameID, _ = strconv.Atoi(item.GameID)
	xmlFeedGame.sport = item.Sport

	//construct hasStarted and isComplete
	xmlFeedGame.hasStarted = xmlFootballSource.constructHasStarted(item.Time)
	xmlFeedGame.isComplete = xmlFootballSource.constructIsComplete(xmlFootballGame)

	//calculate the phase
	phase := xmlFootballSource.calculatePhase(xmlFootballGame, xmlFeedGame.hasStarted, xmlFeedGame.isComplete)

	//construct home and visiting scores
	xmlFeedGame.homeScore = xmlFootballSource.constructHomeScore(xmlFootballGame)
	xmlFeedGame.visitingScore = xmlFootballSource.constructVisitingScore(xmlFootballGame)

	var clock string
	//construct custom data
	xmlFeedGame.customData, clock = xmlFootballSource.constructCustomData(xmlFootballGame, phase)

	//construct regular clock seconds and period for supporting old versions
	xmlFeedGame.period = xmlFootballSource.constructRegularPeriod(phase)
	xmlFeedGame.clockSeconds = xmlFootballSource.constructRegularClock(clock)

	return &xmlFeedGame, nil
}

func (xmlFootballSource *xmlFootballSource) constructRegularPeriod(phase string) int {
	switch phase {
	case "pre":
		return 1
	case "1":
		return 1
	case "2":
		return 2
	case "ht":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "ot":
		return 5
	case "final":
		return 4
	default:
		return -1
	}
}

func (xmlFootballSource *xmlFootballSource) constructRegularClock(clock string) int {
	//we need the seconds but not "00:00" format
	clockArr := strings.Split(clock, ":")
	if clockArr == nil || len(clockArr) != 2 {
		log.Println("xmlfootball: constructRegularClock -> error processing the clock - nil or size != 2")
		return -1
	}

	minutes, err := strconv.Atoi(clockArr[0])
	if err != nil {
		log.Println("xmlfootball: constructRegularClock -> error processing the clock - cannot convert minutes to int " + err.Error())
		return -1
	}
	seconds, err := strconv.Atoi(clockArr[1])
	if err != nil {
		log.Println("xmlfootball: constructRegularClock -> error processing the clock - cannot convert seconds to int " + err.Error())
		return -1
	}
	return (minutes * 60) + seconds
}

func (xmlFootballSource *xmlFootballSource) constructCustomData(xmlData *xmlFootballGame, phase string) (string, string) {

	downtogo := xmlFootballSource.getXMLDowntogo(xmlData)

	hasBall := xmlFootballSource.getPossession(phase, downtogo)
	lastPlay := xmlFootballSource.getLastPlay(phase, downtogo)
	clock := xmlFootballSource.getClock(phase, downtogo)
	phaseLabel := xmlFootballSource.getDisplayPhase(phase)

	footballCustomData := footballCustomData{Possession: hasBall, LastPlay: lastPlay, Clock: clock, Phase: phaseLabel}
	var data []byte
	data, err := json.Marshal(footballCustomData)
	if err != nil {
		log.Printf("xmlfootball: constructCustomData() -> %s\n" + err.Error())
		return "", ""
	}
	value := string(data)
	return value, clock
}

// calculatePhase gives one of the following - pre, 1, 2, ht, 3, 4, ot, final
func (xmlFootballSource *xmlFootballSource) calculatePhase(xmlData *xmlFootballGame, hasStarted bool, isComplete bool) string {
	//check for pre
	if !hasStarted {
		log.Println("xmlfootball: calculatePhase -> pre")
		return "pre"
	}
	//check for final
	if isComplete {
		log.Println("xmlfootball: calculatePhase -> final")
		return "final"
	}

	//the game is started but not completed - we have one of 1, 2, ht, 3, 4, ot

	//get what we have in the xml for quarter
	downtogo := xmlFootballSource.getXMLDowntogo(xmlData)
	if downtogo == nil || len(downtogo.Qtr) == 0 {
		log.Println("xmlfootball: calculatePhase -> still not added a quarter, this means we need to return 1 - first quarter")
		return "1"
	}

	//check for ht
	isHalfTime := xmlFootballSource.isHalfTime(downtogo)
	if isHalfTime {
		log.Println("xmlfootball: calculatePhase -> ht")
		return "ht"
	}

	//check for ot
	isOverTime := xmlFootballSource.isOverTime(downtogo)
	if isOverTime {
		log.Println("xmlfootball: calculatePhase -> ot")
		return "ot"
	}

	//return the quarter - 1, 2, 3 or 4
	log.Printf("xmlfootball: calculatePhase -> return %s", downtogo.Qtr)
	return downtogo.Qtr
}

func (xmlFootballSource *xmlFootballSource) isOverTime(downtogo *xmlFootballDowntogo) bool {
	if downtogo == nil {
		log.Println("xmlfootball: isOverTime -> downtogo is nil")
		return false
	}
	quarter := downtogo.Qtr
	if len(quarter) == 0 {
		log.Println("xmlfootball: isOverTime -> quarter is empty")
		return false
	}
	qrt, err := strconv.Atoi(quarter)
	if err != nil {
		// we do not know the exact format. It could be "5", "ot", "OT" etc, so on parse error we return true
		log.Printf("xmlfootball: isOverTime -> error parsing %s\terror:%s", quarter, err.Error())
		return true
	}
	// the parsing is successfull so we need to check if it is == 5
	return qrt == 5
}

func (xmlFootballSource *xmlFootballSource) isHalfTime(downtogo *xmlFootballDowntogo) bool {
	if downtogo == nil {
		log.Println("xmlfootball: isHalfTime -> downtogo is nil")
		return false
	}
	quarter := downtogo.Qtr
	clock := downtogo.Clock
	return quarter == "2" && clock == "00:00"
}

func (xmlFootballSource *xmlFootballSource) getXMLDowntogo(xmlData *xmlFootballGame) *xmlFootballDowntogo {
	if xmlData == nil {
		log.Printf("xmlfootball: getXMLDowntogo -> the xml data is nil\n")
		return nil
	}
	playsData := xmlData.Plays
	downtogo := playsData.Downtogo
	return &downtogo
}

func (xmlFootballSource *xmlFootballSource) getClock(phase string, downtogo *xmlFootballDowntogo) string {
	if phase == "pre" || phase == "ht" || phase == "final" {
		//do not show for pre, ht and final
		return ""
	}
	if downtogo == nil || len(downtogo.Clock) == 0 {
		//check if the game is started(phase == 1) - there is not added play yet.
		if phase == "1" {
			log.Println("xmlfootball: getClock -> still not added a clock, this means we need to return 15:00")
			return "15:00"
		}
		return ""
	}
	var result string
	clock := downtogo.Clock
	if clock != "00:00" {
		result = clock + " remaining" //we want 02:52 remaining
	} else {
		result = clock
	}
	return result
}

func (xmlFootballSource *xmlFootballSource) getPossession(phase string, downtogo *xmlFootballDowntogo) string {
	if phase == "pre" || phase == "final" {
		//do not show for pre and final
		return ""
	}
	if downtogo == nil {
		return ""
	}
	return downtogo.Hasball
}

func (xmlFootballSource *xmlFootballSource) getLastPlay(phase string, downtogo *xmlFootballDowntogo) string {
	if !xmlFootballSource.config.FootballConfig.LastPlayEnabled {
		//it is disabled
		return ""
	}
	if phase == "pre" || phase == "ht" || phase == "final" {
		//do not show for pre, ht and final
		return ""
	}
	if downtogo == nil {
		return ""
	}
	return downtogo.LastPlay
}

func (xmlFootballSource *xmlFootballSource) getDisplayPhase(phase string) string {
	return xmlFootballSource.config.GetFootballPhaseLabel(phase)
}

func (xmlFootballSource *xmlFootballSource) constructHasStarted(startTime time.Time) bool {
	if startTime.IsZero() {
		log.Println("xmlfootball: constructHasStarted -> start time is zero") //should not happen, return true
		return true
	}

	startInMilliSeconds := startTime.UnixNano() / int64(time.Millisecond)
	nowInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("xmlfootball: constructHasStarted -> start %d\t now %d", startInMilliSeconds, nowInMilliSeconds)

	if nowInMilliSeconds >= startInMilliSeconds {
		return true
	}
	return false
}

func (xmlFootballSource *xmlFootballSource) constructIsComplete(xmlData *xmlFootballGame) bool {
	if xmlData == nil {
		log.Println("xmlfootball: constructIsComplete -> error processing isComplete - the xml data is nil")
		return false
	}
	playsData := xmlData.Plays
	quarters := playsData.Quarters
	if len(quarters) <= 0 {
		log.Println("xmlfootball: constructIsComplete -> quarters is nil or empty")
		return false
	}
	//get the latest quarter
	latestQuarter := quarters[len(quarters)-1]
	scores := latestQuarter.Scores
	if len(scores) <= 0 {
		log.Println("xmlfootball: constructIsComplete -> scores is nil or empty")
		return false
	}
	//get the latest score for the latest quarter
	latestScore := scores[len(scores)-1]
	if latestScore.Final == "Y" {
		return true
	}
	return false
}

func (xmlFootballSource *xmlFootballSource) constructHomeScore(xmlData *xmlFootballGame) int {
	latestScore, err := xmlFootballSource.findLatestScore(xmlData)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	homeScore, err := strconv.Atoi(latestScore.Hscore)
	if err != nil {
		log.Println("xmlfootball: Error processing the home score - cannot convert the score to int " + err.Error())
		return 0
	}
	return homeScore
}

func (xmlFootballSource *xmlFootballSource) constructVisitingScore(xmlData *xmlFootballGame) int {
	latestScore, err := xmlFootballSource.findLatestScore(xmlData)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	visitingScore, err := strconv.Atoi(latestScore.Vscore)
	if err != nil {
		log.Println("xmlfootball: Error processing the visiting score - cannot convert the score to int " + err.Error())
		return 0
	}
	return visitingScore
}

func (xmlFootballSource *xmlFootballSource) findLatestScore(xmlData *xmlFootballGame) (*xmlFootballScore, error) {
	if xmlData == nil {
		return nil, errors.New("xmlfootball: findLatestScore: the xml data is nil")
	}
	scores := xmlData.Scores
	scoresList := scores.ScoresList
	if len(scoresList) <= 0 {
		return nil, errors.New("xmlfootball: findLatestScore: the scores list is nil or empty")
	}
	latest := scoresList[len(scoresList)-1]
	return &latest, nil
}

func (xmlFootballSource *xmlFootballSource) isForGame(xml *xmlFootballGame, item *sidearmModel.LiveGameItem) bool {
	if xml == nil || item == nil {
		log.Println("xmlfootball: isForGame -> xml or item is nil")
		return false
	}
	//date="9/21/2019"
	date := xml.Venue.Date

	dateMonth, dateDay, dateYear, err := parseDate(date)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	itemTime := item.Time
	location, err := time.LoadLocation("America/Chicago")
	if err == nil {
		itemTime = itemTime.In(location)
	}
	itemData := itemTime.Format("01/02/2006")
	itemDateMonth, itemDateDay, itemDateYear, err := parseDate(itemData)

	return (dateMonth == itemDateMonth) && (dateDay == itemDateDay) && (dateYear == itemDateYear)
}

func (xmlFootballSource *xmlFootballSource) printXMLFootballGame(xmlFootballGame *xmlFootballGame) {
	venue := xmlFootballGame.Venue
	log.Printf("venue - date:%s\n", venue.Date)

	xmlFootballSource.printXMLFootballScores(xmlFootballGame.Scores)
	xmlFootballSource.printXMLFootballPlays(xmlFootballGame.Plays)
}

func (xmlFootballSource *xmlFootballSource) printXMLFootballScores(scores xmlFootballScores) {
	scoresList := scores.ScoresList
	if scoresList == nil {
		log.Println("xmlfootball: No scores")
		return
	}

	for _, score := range scoresList {
		log.Printf("score - quarter:%s vscore:%s hscore:%s\n", score.Quarter, score.Vscore, score.Hscore)
	}
}

func (xmlFootballSource *xmlFootballSource) printXMLFootballPlays(data xmlFootballPlays) {
	downtogo := data.Downtogo
	log.Printf("downtogo - hasball:%s quarter:%s clock:%s lastplay:%s\n", downtogo.Hasball, downtogo.Qtr, downtogo.Clock, downtogo.LastPlay)

	quarters := data.Quarters
	if quarters == nil {
		log.Println("\tNo quarters")
		return
	}

	for _, quarter := range quarters {
		qtrNumber := quarter.Number
		log.Println("Quarter:" + qtrNumber)

		plays := quarter.Plays
		if plays == nil {
			log.Println("\tNo playes")
		} else {
			for _, play := range plays {
				log.Printf("\tplay - score:%s vscore:%s hscore:%s clock:%s\n", play.Score, play.Vscore, play.Hscore, play.Clock)
			}
		}

		scores := quarter.Scores
		if scores == nil {
			log.Println("\tNo scores")
		} else {
			for _, score := range scores {
				log.Printf("\tscore - final:%s v:%s h:%s\n", score.Final, score.V, score.H)
			}
		}
	}
}
