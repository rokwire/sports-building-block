package source

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"sport/core/model"
	sidearmModel "sport/driven/provider/sidearm/model"
	"strconv"
	"time"
)

type xmlFeedVolleyballGame struct {
	gameID        int
	sport         string
	clockSeconds  int
	period        int
	hasStarted    bool
	isComplete    bool
	homeScore     int
	visitingScore int
	customData    string
}

func (game *xmlFeedVolleyballGame) GetType() string {
	return game.sport
}

func (game *xmlFeedVolleyballGame) GetGameID() int {
	return game.gameID
}

func (game *xmlFeedVolleyballGame) GetPath() string {
	return game.sport
}

func (game *xmlFeedVolleyballGame) GetHasStarted() bool {
	return game.hasStarted
}

func (game *xmlFeedVolleyballGame) GetIsComplete() bool {
	return game.isComplete
}

func (game *xmlFeedVolleyballGame) GetClockSeconds() int {
	return game.clockSeconds
}

func (game *xmlFeedVolleyballGame) GetPeriod() int {
	return game.period
}

func (game *xmlFeedVolleyballGame) GetHomeScore() int {
	return game.homeScore
}

func (game *xmlFeedVolleyballGame) GetVisitingScore() int {
	return game.visitingScore
}

func (game *xmlFeedVolleyballGame) GetCustomData() string {
	return game.customData
}

func (game *xmlFeedVolleyballGame) Encode() map[string]string {

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

type xmlVolleyballGame struct {
	XMLName   xml.Name             `xml:"vbgame"`
	Generated string               `xml:"generated,attr"`
	Venue     xmlVolleyballVenue   `xml:"venue"`
	Status    *xmlVolleyballStatus `xml:"status"`
}

type xmlVolleyballVenue struct {
	XMLName xml.Name `xml:"venue"`
	Date    string   `xml:"date,attr"`
}

type xmlVolleyballStatus struct {
	XMLName  xml.Name `xml:"status"`
	Complete string   `xml:"complete,attr"`
	VSCore   string   `xml:"vscore,attr"`
	HScore   string   `xml:"hscore,attr"`
	Game     string   `xml:"game,attr"`
	Serving  string   `xml:"serving,attr"`
	VPoints  string   `xml:"vpoints,attr"`
	HPoints  string   `xml:"hpoints,attr"`
}

type volleyballCustomData struct {
	HasExtraData bool
	Phase        string
	PhaseLabel   string
	HScore       string
	VScore       string
	HPoints      string
	VPoints      string
	Serving      string
}

type xmlVolleyballSource struct {
	config  Config
	ftpConn ftpConn
}

func (xmlVolleyballSource *xmlVolleyballSource) updateConfig(config Config) {
	log.Println("xmlvolleyball: UpdateConfig -> config updated in xml volleyball source")
	xmlVolleyballSource.config = config
}

func (xmlVolleyballSource *xmlVolleyballSource) load(item *sidearmModel.LiveGameItem) (model.LiveGame, error) {
	//1. load the xml data
	xmlData, err := xmlVolleyballSource.ftpConn.load("/Rokwire_FTP/wvball")
	if err != nil {
		return nil, err
	}

	//2. unmarshal
	var xmlVolleyballGame *xmlVolleyballGame
	err = xml.Unmarshal(xmlData, &xmlVolleyballGame)
	if err != nil {
		return nil, err
	}

	//xmlVolleyballSource.printXMLVolleyballGame(xmlVolleyballGame)

	//3. check if the dowloaded xml is for this game. The only way we could check is to compare the dates!
	if xmlVolleyballSource.config.GetVolleyballDateCheck() && !xmlVolleyballSource.isForGame(xmlVolleyballGame, item) {
		return nil, errors.New("xmlvolleyball loadFromXML -> the xml is not for this game")
	}

	//4. construct xmlFeedFootballGame
	var xmlFeedGame xmlFeedVolleyballGame
	xmlFeedGame.gameID, _ = strconv.Atoi(item.GameID)
	xmlFeedGame.sport = item.Sport

	//construct hasStarted and isComplete
	xmlFeedGame.hasStarted = xmlVolleyballSource.constructHasStarted(item.Time)
	xmlFeedGame.isComplete = xmlVolleyballSource.constructIsComplete(item.Time, xmlVolleyballGame)

	//construct period
	xmlFeedGame.period = xmlVolleyballSource.constructPeriod(xmlVolleyballGame)

	//construct home and visiting score
	xmlFeedGame.homeScore = xmlVolleyballSource.constructScore("H", xmlVolleyballGame)
	xmlFeedGame.visitingScore = xmlVolleyballSource.constructScore("V", xmlVolleyballGame)

	//construct custom data
	xmlFeedGame.customData = xmlVolleyballSource.constructCustomData(xmlVolleyballGame, xmlFeedGame.hasStarted, xmlFeedGame.isComplete)

	return &xmlFeedGame, nil
}

func (xmlVolleyballSource *xmlVolleyballSource) constructCustomData(xmlData *xmlVolleyballGame, started bool, completed bool) string {
	volleyballCustomData := xmlVolleyballSource.createExtraData(xmlData, started, completed)
	var data []byte
	data, err := json.Marshal(volleyballCustomData)
	if err != nil {
		log.Printf("xmlvolleyball constructCustomData() -> %s\n" + err.Error())
		return ""
	}
	value := string(data)
	return value
}

func (xmlVolleyballSource *xmlVolleyballSource) createExtraData(xmlData *xmlVolleyballGame, started bool, completed bool) volleyballCustomData {
	if xmlData == nil {
		log.Println("xmlvolleyball createExtraData -> xml is nil")
		return volleyballCustomData{HasExtraData: false}
	}
	status := xmlData.Status
	if status == nil {
		log.Println("xmlvolleyball createExtraData -> status is nil")
		return volleyballCustomData{HasExtraData: false}
	}

	phase, phaseLabel := xmlVolleyballSource.calculatePhase(status, started, completed)
	hScore := status.HScore
	vScore := status.VSCore
	hPoints := xmlVolleyballSource.getPoints("H", phase, status)
	vPoints := xmlVolleyballSource.getPoints("V", phase, status)
	serving := xmlVolleyballSource.getServing(phase, status)
	volleyballCustomData := volleyballCustomData{HasExtraData: true, Phase: phase, PhaseLabel: phaseLabel,
		HScore: hScore, VScore: vScore, HPoints: hPoints, VPoints: vPoints, Serving: serving}
	return volleyballCustomData
}

func (xmlVolleyballSource *xmlVolleyballSource) getPoints(vh string, phase string, status *xmlVolleyballStatus) string {
	if status == nil {
		log.Println("xmlvolleyball getPoints -> status is nil")
		return ""
	}

	//send points only fi the game is active
	if phase == "pre" || phase == "final" {
		log.Println("xmlvolleyball getPoints -> do not send points for pre or final phases")
		return ""
	}

	switch vh {
	case "H":
		return status.HPoints
	case "V":
		return status.VPoints
	default:
		return ""
	}
}

func (xmlVolleyballSource *xmlVolleyballSource) getServing(phase string, status *xmlVolleyballStatus) string {
	if status == nil {
		log.Println("xmlvolleyball getServing -> status is nil")
		return ""
	}

	//send points only fi the game is active
	if phase == "pre" || phase == "final" {
		log.Println("xmlvolleyball getServing -> do not send serving for pre or final phases")
		return ""
	}

	return status.Serving
}

func (xmlVolleyballSource *xmlVolleyballSource) calculatePhase(status *xmlVolleyballStatus, started bool, completed bool) (string, string) {
	config := xmlVolleyballSource.config
	//check for pre
	if !started {
		log.Println("xmlvolleyball calculatePhase -> pre")
		return "pre", config.GetVolleyballPhaseLabel("pre")
	}
	//check for final
	if completed {
		log.Println("xmlvolleyball calculatePhase -> final")
		return "final", config.GetVolleyballPhaseLabel("final")
	}

	//the game is started but not completed - we have during game phases - 1, 2, 3..

	//check if we have data
	if status == nil {
		log.Println("xmlvolleyball calculatePhase -> for some reasons we do not have data yet, so return pre")
		return "pre", config.GetVolleyballPhaseLabel("pre")
	}

	game := status.Game

	//check if the game is number
	isNumber := isNumber(game)
	if !isNumber {
		log.Println("xmlvolleyball calculatePhase -> for some reasons the game is not number, so return pre")
		return "pre", config.GetVolleyballPhaseLabel("pre")
	}

	phaseLabel := getOrdinal(game) + " " + config.GetVolleyballPhaseLabel("game_name")
	return game, phaseLabel
}

func (xmlVolleyballSource *xmlVolleyballSource) constructScore(vh string, xmlData *xmlVolleyballGame) int {
	if xmlData == nil {
		log.Println("xmlvolleyball constructScore -> xml is nil")
		return 0
	}
	status := xmlData.Status
	if status == nil {
		log.Println("xmlvolleyball constructScore -> status is nil")
		return 0
	}
	var score string
	if vh == "H" {
		score = status.HScore
	} else {
		score = status.VSCore
	}

	result, err := strconv.Atoi(score)
	if err != nil {
		log.Println("xmlvolleyball - cannot convert the score to int " + err.Error())
		return 0
	}

	return result
}

func (xmlVolleyballSource *xmlVolleyballSource) constructPeriod(xmlData *xmlVolleyballGame) int {
	if xmlData == nil {
		log.Println("xmlvolleyball constructPeriod -> xml is nil, so return 1")
		return 1
	}
	status := xmlData.Status
	if status == nil {
		log.Println("xmlvolleyball constructPeriod -> status is nil")
		return 1
	}
	game, err := strconv.Atoi(status.Game)
	if err != nil {
		log.Println("xmlvolleyball - cannot convert the period to int " + err.Error())
		return 1
	}
	return game
}

func (xmlVolleyballSource *xmlVolleyballSource) constructHasStarted(startTime time.Time) bool {
	if startTime.IsZero() {
		log.Println("xmlvolleyball constructHasStarted -> start time is zero") //should not happen, return true
		return true
	}

	startInMilliSeconds := startTime.UnixNano() / int64(time.Millisecond)
	nowInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("xmlvolleyball constructHasStarted -> start %d\t now %d", startInMilliSeconds, nowInMilliSeconds)

	if nowInMilliSeconds >= startInMilliSeconds {
		return true
	}
	return false
}

func (xmlVolleyballSource *xmlVolleyballSource) constructIsComplete(startTime time.Time, xmlData *xmlVolleyballGame) bool {
	if xmlData == nil {
		log.Println("xmlvolleyball constructIsComplete -> error processing isComplete - the xml data is nil")
		return false
	}
	status := xmlData.Status
	if status == nil {
		log.Println("xmlvolleyball constructIsComplete -> status is nil")
		return false
	}
	complete := status.Complete
	if complete == "Y" {
		return true
	}
	//for some reasons sometimes it does not mark the complete status with "Y", so we need to check by time
	afterDay := startTime.Add(time.Hour * time.Duration(24))
	afterDayInMilliSeconds := afterDay.UnixNano() / int64(time.Millisecond)
	nowInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	if nowInMilliSeconds >= afterDayInMilliSeconds {
		log.Println("xmlvolleyball constructIsComplete -> mark as complete because it is a day later")
		return true
	}

	return false
}

func (xmlVolleyballSource *xmlVolleyballSource) isForGame(xml *xmlVolleyballGame, item *sidearmModel.LiveGameItem) bool {
	if xml == nil || item == nil {
		log.Println("isForGame -> xml or item is nil")
		return false
	}
	//date="11/01/2019"
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

func (xmlVolleyballSource *xmlVolleyballSource) printXMLVolleyballGame(xmlVolleyballGame *xmlVolleyballGame) {
	venue := xmlVolleyballGame.Venue
	log.Printf("venue - date:%s\n", venue.Date)

	status := xmlVolleyballGame.Status
	log.Printf("status - complete:%s vscore:%s hscore:%s game:%s serving:%s vpoints:%s hpoints:%s\n",
		status.Complete, status.VSCore, status.HScore, status.Game, status.Serving, status.VPoints, status.HPoints)
}

func newXMLVolleyballSource(config Config, ftpHost string, ftpUser string, ftpPassword string) xmlVolleyballSource {
	var xmlVolleyballSource xmlVolleyballSource
	xmlVolleyballSource.config = config
	xmlVolleyballSource.ftpConn = newFTPConn(ftpHost, ftpUser, ftpPassword)
	return xmlVolleyballSource
}
