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

type xmlFeedBasketballGame struct {
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

func (game *xmlFeedBasketballGame) GetType() string {
	return game.sport
}

func (game *xmlFeedBasketballGame) GetGameID() int {
	return game.gameID
}

func (game *xmlFeedBasketballGame) GetPath() string {
	return game.sport
}

func (game *xmlFeedBasketballGame) GetHasStarted() bool {
	return game.hasStarted
}

func (game *xmlFeedBasketballGame) GetIsComplete() bool {
	return game.isComplete
}

func (game *xmlFeedBasketballGame) GetClockSeconds() int {
	return game.clockSeconds
}

func (game *xmlFeedBasketballGame) GetPeriod() int {
	return game.period
}

func (game *xmlFeedBasketballGame) GetHomeScore() int {
	return game.homeScore
}

func (game *xmlFeedBasketballGame) GetVisitingScore() int {
	return game.visitingScore
}

func (game *xmlFeedBasketballGame) GetCustomData() string {
	return game.customData
}

func (game *xmlFeedBasketballGame) Encode() map[string]string {

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

type xmlBasketballGame struct {
	XMLName   xml.Name            `xml:"bbgame"`
	Generated string              `xml:"generated,attr"`
	Venue     xmlBasketballVenue  `xml:"venue"`
	Status    xmlBasketballStatus `xml:"status"`
	Teams     []xmlBasketballTeam `xml:"team"`
	Plays     xmlBasketballPlays  `xml:"plays"`
}

type xmlBasketballVenue struct {
	XMLName xml.Name           `xml:"venue"`
	Date    string             `xml:"date,attr"`
	Start   string             `xml:"start,attr"`
	Rules   xmlBasketballRules `xml:"rules"`
}

type xmlBasketballRules struct {
	XMLName xml.Name `xml:"rules"`
	Prds    string   `xml:"prds,attr"`
}

type xmlBasketballStatus struct {
	XMLName  xml.Name `xml:"status"`
	Complete string   `xml:"complete,attr"`
	Period   string   `xml:"period,attr"`
	Clock    string   `xml:"clock,attr"`
	Running  string   `xml:"running,attr"`
}

type xmlBasketballTeam struct {
	XMLName   xml.Name               `xml:"team"`
	VH        string                 `xml:"vh,attr"`
	Name      string                 `xml:"name,attr"`
	Linescore xmlBasketballLinescore `xml:"linescore"`
	Players   []xmlBasketballPlayer  `xml:"player"`
}

type xmlBasketballPlayer struct {
	XMLName   xml.Name `xml:"player"`
	Name      string   `xml:"name,attr"`
	Checkname string   `xml:"checkname,attr"`
}

type xmlBasketballLinescore struct {
	XMLName xml.Name `xml:"linescore"`
	Score   string   `xml:"score,attr"`
}

type xmlBasketballSource struct {
	config  Config
	ftpConn ftpConn
}

type xmlBasketballPlays struct {
	XMLName xml.Name              `xml:"plays"`
	Periods []xmlBasketballPeriod `xml:"period"`
}

type xmlBasketballPeriod struct {
	XMLName xml.Name            `xml:"period"`
	Number  string              `xml:"number,attr"`
	Plays   []xmlBasketballPlay `xml:"play"`
}

type basketballCustomData struct {
	Clock    string
	Phase    string
	LastPlay string
}

type xmlBasketballPlay struct {
	XMLName   xml.Name `xml:"play"`
	VH        string   `xml:"vh,attr"`
	Time      string   `xml:"time,attr"`
	UNI       string   `xml:"uni,attr"`
	Team      string   `xml:"team,attr"`
	Checkname string   `xml:"checkname,attr"`
	Action    string   `xml:"action,attr"`
	Type      string   `xml:"type,attr"`
	Side      string   `xml:"side,attr"`
}

func newXMLBasketballSource(config Config, ftpHost string, ftpUser string, ftpPassword string) xmlBasketballSource {
	var xmlBasketballSource xmlBasketballSource
	xmlBasketballSource.config = config
	xmlBasketballSource.ftpConn = newFTPConn(ftpHost, ftpUser, ftpPassword)
	return xmlBasketballSource
}

func (xmlBasketballSource *xmlBasketballSource) updateConfig(config Config) {
	log.Println("xmlbasketball: UpdateConfig -> config updated in xml basketball source")
	xmlBasketballSource.config = config
}

func (xmlBasketballSource *xmlBasketballSource) load(item *sidearmModel.LiveGameItem) (model.LiveGame, error) {
	//1. load the xml data
	xmlData, err := xmlBasketballSource.ftpConn.load("/Rokwire_FTP/" + item.Sport) // mbball or wbball
	if err != nil {
		return nil, err
	}

	//2. unmarshal
	var xmlBasketballGame *xmlBasketballGame
	err = xml.Unmarshal(xmlData, &xmlBasketballGame)
	if err != nil {
		return nil, err
	}

	//xmlBasketballSource.printXMLFootballGame(xmlBasketballGame)

	//3. check if the dowloaded xml is for this game. The only way we could check is to compare the dates!
	if xmlBasketballSource.config.GetBasketballDateCheck(item.Sport) && !xmlBasketballSource.isForGame(xmlBasketballGame, item) {
		return nil, errors.New("xmlbasketball loadFromXML -> the xml is not for this game")
	}

	//4. construct xmlFeedFootballGame
	var xmlFeedGame xmlFeedBasketballGame
	xmlFeedGame.gameID, _ = strconv.Atoi(item.GameID)
	xmlFeedGame.sport = item.Sport

	//construct hasStarted and isComplete
	xmlFeedGame.hasStarted = xmlBasketballSource.constructHasStarted(item.Time)
	xmlFeedGame.isComplete = xmlBasketballSource.constructIsComplete(item.Time, xmlBasketballGame)

	//calculate the phase
	phase := xmlBasketballSource.calculatePhase(xmlBasketballGame, xmlFeedGame.hasStarted, xmlFeedGame.isComplete)

	//construct the regular period and clock - old verisons
	xmlFeedGame.period = xmlBasketballSource.constructRegularPeriod(xmlBasketballGame.Status.Period)
	xmlFeedGame.clockSeconds = xmlBasketballSource.constructRegularClock(xmlBasketballGame.Status.Clock)

	//construct home and visiting scores
	xmlFeedGame.homeScore = xmlBasketballSource.constructScore("H", xmlBasketballGame)
	xmlFeedGame.visitingScore = xmlBasketballSource.constructScore("V", xmlBasketballGame)

	//construct custom data
	xmlFeedGame.customData = xmlBasketballSource.constructCustomData(xmlBasketballGame, phase, item.Sport)

	return &xmlFeedGame, nil
}

func (xmlBasketballSource *xmlBasketballSource) isForGame(xml *xmlBasketballGame, item *sidearmModel.LiveGameItem) bool {
	if xml == nil || item == nil {
		log.Println("xmlbasketball isForGame -> xml or item is nil")
		return false
	}
	//date="9/21/2019"
	date := xml.Venue.Date // it is UTC

	dateMonth, dateDay, dateYear, err := parseDate(date)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	itemTime := item.Time
	itemData := itemTime.Format("01/02/2006")
	itemDateMonth, itemDateDay, itemDateYear, err := parseDate(itemData)

	return (dateMonth == itemDateMonth) && (dateDay == itemDateDay) && (dateYear == itemDateYear)
}

func (xmlBasketballSource *xmlBasketballSource) constructCustomData(xmlData *xmlBasketballGame, phase string, sport string) string {
	clock := xmlBasketballSource.getClock(xmlData, phase)
	phaseLabel := xmlBasketballSource.getDisplayPhase(phase, sport)
	lastPlay := xmlBasketballSource.getLastPlay(xmlData, phase, sport)

	basketballCustomData := basketballCustomData{Clock: clock, Phase: phaseLabel, LastPlay: lastPlay}
	var data []byte
	data, err := json.Marshal(basketballCustomData)
	if err != nil {
		log.Printf("xmlbasketball constructCustomData() -> %s\n" + err.Error())
		return ""
	}
	value := string(data)
	return value
}

func (xmlBasketballSource *xmlBasketballSource) getLastPlay(xmlData *xmlBasketballGame, phase string, sport string) string {
	if !xmlBasketballSource.config.GetBasketballLastPlay(sport) {
		//it is disabled
		return ""
	}
	if phase == "pre" || phase == "final" {
		//do not show for pre and final
		return ""
	}
	lastPlay := xmlBasketballSource.findLatestPlay(xmlData)
	if lastPlay == nil {
		//there is no plays yet, so return empty string
		return ""
	}
	return xmlBasketballSource.formatPlay(xmlData, lastPlay)
}

func (xmlBasketballSource *xmlBasketballSource) formatPlay(xmlData *xmlBasketballGame, play *xmlBasketballPlay) string {
	if play == nil {
		return ""
	}
	var result string
	if len(play.Team) > 0 {
		result = result + play.Team + " team,"
	}
	if len(play.Checkname) > 0 && play.Checkname != "TEAM" {
		player := xmlBasketballSource.findPlayer(play.VH, play.Checkname, xmlData)
		var playerLabel string
		if player != nil {
			playerLabel = player.Name
		} else {
			playerLabel = play.Checkname
		}
		result = result + " " + playerLabel + ","
	}
	if len(play.Action) > 0 {
		result = result + " action - " + play.Action + ","
	}
	if len(play.Type) > 0 {
		result = result + " type - " + play.Type + ","
	}
	if len(play.Side) > 0 {
		result = result + " side - " + play.Side + ","
	}

	if len(result) > 1 {
		result = result[0 : len(result)-1] //remove the last ","
	}
	return result
}

func (xmlBasketballSource *xmlBasketballSource) findPlayer(vh string, checkname string, xmlData *xmlBasketballGame) *xmlBasketballPlayer {
	team := xmlBasketballSource.findTeam(vh, xmlData)
	if team == nil {
		log.Println("xmlbasketball findPlayer -> no team")
		return nil
	}
	players := team.Players
	if players == nil {
		log.Println("xmlbasketball findPlayer -> no players")
		return nil
	}
	for _, player := range players {
		if player.Checkname == checkname {
			return &player
		}
	}
	log.Printf("xmlbasketball findPlayer -> no player for checkname %s\n", checkname)
	return nil
}

func (xmlBasketballSource *xmlBasketballSource) findLatestPlay(xmlData *xmlBasketballGame) *xmlBasketballPlay {
	if xmlData == nil {
		log.Println("xmlbasketball findLatestPlay -> error processing findLatestPlay - the xml data is nil")
		return nil
	}
	periods := xmlData.Plays.Periods
	if periods == nil {
		log.Println("xmlbasketball findLatestPlay -> there is no periods yet")
		return nil
	}
	lastPeriod := periods[len(periods)-1] // get the latest period
	plays := lastPeriod.Plays
	if plays == nil {
		log.Println("xmlbasketball findLatestPlay -> there is no plays yet")
		return nil
	}
	lastPlay := plays[len(plays)-1] // get the latest play
	return &lastPlay
}

func (xmlBasketballSource *xmlBasketballSource) getClock(xmlData *xmlBasketballGame, phase string) string {
	if phase == "pre" || phase == "final" {
		//do not show for pre and final
		return ""
	}
	if xmlData == nil {
		log.Println("xmlbasketball getClock -> error processing getClock - the xml data is nil")
		return ""
	}
	var result string
	clock := xmlData.Status.Clock
	if clock != "00:00" {
		result = clock + " remaining" //we want 02:52 remaining
	} else {
		result = clock
	}
	return result
}

func (xmlBasketballSource *xmlBasketballSource) getDisplayPhase(phase string, sport string) string {
	switch sport {
	case "mbball":
		return xmlBasketballSource.config.GetMBasketballPhaseLabel(phase)
	case "wbball":
		return xmlBasketballSource.config.GetWBasketballPhaseLabel(phase)
	default:
		return ""
	}
}

func (xmlBasketballSource *xmlBasketballSource) constructScore(vh string, xmlData *xmlBasketballGame) int {
	team := xmlBasketballSource.findTeam(vh, xmlData)
	if team == nil {
		log.Println("xmlbasketball constructScore -> team is nil")
		return 0
	}
	score := team.Linescore.Score
	iScore, err := strconv.Atoi(score)
	if err != nil {
		log.Printf("xmlbasketball isOverTime -> score cannot be casted %s", score)
		return 0
	}
	return iScore
}

func (xmlBasketballSource *xmlBasketballSource) findTeam(vh string, xmlData *xmlBasketballGame) *xmlBasketballTeam {
	if xmlData == nil {
		log.Println("xmlbasketball findTeam -> error processing findTeam - the xml data is nil")
		return nil
	}
	teams := xmlData.Teams
	if teams == nil {
		log.Println("xmlbasketball findTeam -> teams is nil")
		return nil
	}
	for _, team := range teams {
		if vh == team.VH {
			return &team
		}
	}
	return nil
}

func (xmlBasketballSource *xmlBasketballSource) constructRegularPeriod(period string) int {
	if len(period) == 0 {
		log.Println("xmlbasketball constructRegularPeriod -> period is empty so return 1")
		return 1
	}
	cPeriod, err := strconv.Atoi(period)
	if err != nil {
		log.Println("xmlbasketball constructRegularPeriod -> period cannot be casted so return 1")
		return 1
	}
	return cPeriod
}

func (xmlBasketballSource *xmlBasketballSource) constructRegularClock(clock string) int {
	//we need the seconds but not "00:00" format
	clockArr := strings.Split(clock, ":")
	if clockArr == nil || len(clockArr) != 2 {
		log.Println("xmlbasketball constructRegularClock -> error processing the clock - nil or size != 2")
		return -1
	}

	minutes, err := strconv.Atoi(clockArr[0])
	if err != nil {
		log.Println("xmlbasketball constructRegularClock -> error processing the clock - cannot convert minutes to int " + err.Error())
		return -1
	}
	seconds, err := strconv.Atoi(clockArr[1])
	if err != nil {
		log.Println("xmlbasketball constructRegularClock -> error processing the clock - cannot convert seconds to int " + err.Error())
		return -1
	}
	return (minutes * 60) + seconds
}

func (xmlBasketballSource *xmlBasketballSource) constructHasStarted(startTime time.Time) bool {
	if startTime.IsZero() {
		log.Println("xmlbasketball constructHasStarted -> start time is zero") //should not happen, return true
		return true
	}

	startInMilliSeconds := startTime.UnixNano() / int64(time.Millisecond)
	nowInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("xmlbasketball constructHasStarted -> start %d\t now %d", startInMilliSeconds, nowInMilliSeconds)

	if nowInMilliSeconds >= startInMilliSeconds {
		return true
	}
	return false
}

func (xmlBasketballSource *xmlBasketballSource) constructIsComplete(startTime time.Time, xmlData *xmlBasketballGame) bool {
	if xmlData == nil {
		log.Println("xmlbasketball constructIsComplete -> error processing isComplete - the xml data is nil")
		return false
	}
	status := xmlData.Status
	complete := status.Complete
	if complete == "Y" {
		return true
	}
	//for some reasons sometimes it does not mark the complete status with "Y", so we need to check by time
	afterDay := startTime.Add(time.Hour * time.Duration(24))
	afterDayInMilliSeconds := afterDay.UnixNano() / int64(time.Millisecond)
	nowInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	if nowInMilliSeconds >= afterDayInMilliSeconds {
		log.Println("xmlbasketball constructIsComplete -> mark as complete because it is a day later")
		return true
	}

	return false
}

// calculatePhase gives one of the following:
// men's basketball - pre, 1, 2, ot, final
// women's basketball - pre, 1, 2, 3, 4, ot, final
func (xmlBasketballSource *xmlBasketballSource) calculatePhase(xmlData *xmlBasketballGame, hasStarted bool, isComplete bool) string {
	if xmlData == nil {
		log.Println("xmlbasketball calculatePhase -> error processing calculatePhase - the xml data is nil")
		return "pre"
	}

	//check for pre
	if !hasStarted {
		log.Println("xmlbasketball calculatePhase -> pre")
		return "pre"
	}
	//check for final
	if isComplete {
		log.Println("xmlbasketball calculatePhase -> final")
		return "final"
	}

	//the game is started but not completed - we have one of (1, 2, ot) or one of (1, 2, 3, 4 ot)
	status := xmlData.Status
	period := status.Period

	//check is over time
	isOverTime := xmlBasketballSource.isOverTime(period, xmlData.Venue.Rules.Prds)
	if isOverTime {
		log.Println("xmlbasketball calculatePhase -> ot")
		return "ot"
	}

	log.Printf("xmlbasketball calculatePhase -> return %s", period)
	return period
}

func (xmlBasketballSource *xmlBasketballSource) isOverTime(currentPeriod string, rulesPeriods string) bool {
	if len(currentPeriod) == 0 || len(rulesPeriods) == 0 {
		log.Printf("xmlbasketball isOverTime -> currentPeriod or rulesPeriod is empty - %s %s", currentPeriod, rulesPeriods)
		return false
	}
	cPeriod, err := strconv.Atoi(currentPeriod)
	if err != nil {
		log.Println("xmlbasketball isOverTime -> current period cannot be casted")
		return false
	}
	rPeriod, err := strconv.Atoi(rulesPeriods)
	if err != nil {
		log.Println("xmlbasketball isOverTime -> rules period cannot be casted")
		return false
	}
	return cPeriod > rPeriod
}

func (xmlBasketballSource *xmlBasketballSource) printXMLFootballGame(xmlBasketballGame *xmlBasketballGame) {
	venue := xmlBasketballGame.Venue
	log.Printf("venue - date:%s start:%s\n", venue.Date, venue.Start)

	status := xmlBasketballGame.Status
	log.Printf("status - complete:%s period:%s clock:%s running:%s\n", status.Complete, status.Period, status.Clock, status.Running)

	// teams
	teams := xmlBasketballGame.Teams
	if teams != nil {
		for _, team := range teams {
			log.Printf("\tteam - vh:%s name:%s score:%s\n", team.VH, team.Name, team.Linescore.Score)
		}
	} else {
		log.Println("teams - nil")
	}

	// plays
	plays := xmlBasketballGame.Plays
	periods := plays.Periods
	if periods != nil {
		for _, period := range periods {
			log.Printf("\tperiod - number:%s\n", period.Number)
			periodPlays := period.Plays
			if periodPlays != nil {
				for _, play := range periodPlays {
					log.Printf("\t\tplay - VH:%s Time:%s UNI:%s Team:%s Checkname:%s Time:%s Action:%s Type:%s\n",
						play.VH, play.Time, play.UNI, play.Team, play.Checkname, play.Time, play.Action, play.Type)
				}
			} else {
				log.Println("\t\tno play for this period")
			}
		}
	} else {
		log.Println("periods - nil")
	}
}
