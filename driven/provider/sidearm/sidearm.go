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

package sidearm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sport/core/model"
	"sport/driven/notifications"
	"sport/driven/provider/sidearm/livestats"
	"sport/driven/provider/sidearm/livestats/source"
	sidearmModel "sport/driven/provider/sidearm/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

const host string = "https://fightingillini.com"
const illinoisTeamName string = "Illinois"

// Provider implements Provider interface
type Provider struct {
	mu            sync.Mutex
	stats         livestats.LiveStats
	config        source.Config
	notifications notifications.Notifications
	nextGame      sidearmModel.LiveGameItem
	startedGames  []*sidearmModel.LiveGameItem
	cachedGames   []sidearmModel.Game
	cachedNews    []model.News
}

// NewProvider creates new provider instance
func NewProvider(internalAPIKey string, host string, ftpHost string, ftpUser string, ftpPassword string, appID string, orgID string) *Provider {
	config := source.NewConfig()
	notifications := notifications.New(internalAPIKey, host, appID, orgID)
	stats := livestats.New(notifications, config, ftpHost, ftpUser, ftpPassword, illinoisTeamName)
	return &Provider{stats: stats, config: config, notifications: notifications}
}

// Start Provider
func (p *Provider) Start() {
	go p.processCachedGames()
	go p.processLiveStats()
	p.loadCachedNews()
	go p.processCachedNews()
}

// GetNews retrieves the news from sidearm service
func (p *Provider) GetNews(id *string, sports []string, limit int) ([]model.News, error) {
	return p.loadNews(id, sports, limit)
}

// GetCoaches retrieves the coaches from sidearm service
func (p *Provider) GetCoaches(sport string) ([]model.Coach, error) {
	coachesEndpoint := "/services/coaches_xml.aspx?format=json"

	if sport != "" {
		coachesEndpoint += fmt.Sprintf("&path=%s", sport)
	}

	url := host + coachesEndpoint
	bodyBytes, err := request(http.MethodGet, url, nil)

	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> GetCoaches: Failed to request coaches. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var r sidearmModel.Rosters
	em := json.Unmarshal(bodyBytes, &r)
	if em != nil {
		log.Printf("sidearm -> GetCoaches: Failed to unmarshal response json. Reason: %s", em.Error())
		return nil, em
	}
	var coaches []model.Coach
	rosters := r.Rosters
	if (rosters != nil) && (len(rosters) > 0) {
		for i := 0; i < len(rosters); i++ {
			r := rosters[i]
			photos := buildRosterPhotos(r.Photos)
			coaches = append(coaches, model.Coach{ID: strconv.Itoa(r.StaffID), Name: r.Name, FirstName: r.FirstName, LastName: r.LastName, Email: r.StaffInfo.Email, Phone: r.StaffInfo.Phone, Title: r.StaffInfo.Title, Bio: r.Bio, Photos: photos})
		}
	}

	return coaches, nil
}

// GetPlayers retrieves the players from sidearm service
func (p *Provider) GetPlayers(sport string) ([]model.Player, error) {
	rosterEndpoint := "/services/roster_xml.aspx?format=json"

	if sport != "" {
		rosterEndpoint += fmt.Sprintf("&path=%s", sport)
	}

	url := host + rosterEndpoint
	bodyBytes, err := request(http.MethodGet, url, nil)

	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> GetPlayers: Failed to request players. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var r sidearmModel.Rosters
	em := json.Unmarshal(bodyBytes, &r)
	if em != nil {
		log.Printf("sidearm -> GetPlayers: Failed to unmarshal response json. Reason: %s", em.Error())
		return nil, em
	}
	var players []model.Player
	rosters := r.Rosters
	if (rosters != nil) && (len(rosters) > 0) {
		for i := 0; i < len(rosters); i++ {
			r := rosters[i]
			photos := buildRosterPhotos(r.Photos)
			var captain bool
			if r.PlayerInfo.Captain == "True" {
				captain = true
			} else {
				captain = false
			}
			player := model.Player{ID: r.PlayerID, Name: r.Name, FirstName: r.FirstName, LastName: r.LastName, Uni: r.PlayerInfo.Uni, PosShort: r.PlayerInfo.PosShort, Height: r.PlayerInfo.Height, Wight: r.PlayerInfo.Weight, Gender: r.PlayerInfo.Gender, YearLong: r.PlayerInfo.YearLong, HomeTown: r.PlayerInfo.HomeTown, HighSchool: r.PlayerInfo.HighSchool, Captain: captain, Bio: r.Bio}
			if photos != nil {
				player.Photos = photos
			}
			players = append(players, player)
		}
	}

	return players, nil
}

// GetSocialNetworks retrieves social accounts from sidearm service
func (p *Provider) GetSocialNetworks() ([]model.SportSocial, error) {
	url := host + "/api/assets?operation=sports"
	bodyBytes, err := request(http.MethodGet, url, nil)

	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> GetSocialNetworks: Failed to request social networks. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var s sidearmModel.SportsSocNet
	es := json.Unmarshal(bodyBytes, &s)
	if es != nil {
		log.Printf("sidearm -> GetSocialNetworks: Failed to unmarshal response json. Reason: %s", es.Error())
		return nil, es
	}

	var socNetList []model.SportSocial
	srcSocNet := s.SportsSocial
	if (srcSocNet != nil) && (len(srcSocNet) > 0) {
		for i := 0; i < len(srcSocNet); i++ {
			r := srcSocNet[i]
			socNet := model.SportSocial{SportShortName: r.ShortName, TwitterName: r.TwitterName, InstagramName: r.InstagramName, FacebookPage: r.FacebookPage}
			socNetList = append(socNetList, socNet)
		}
	}

	return socNetList, nil
}

// GetGames retrieves games from sidearm
func (p *Provider) GetGames(sports []string, id *string, startDate *string, endDate *string, limit int) ([]model.Game, error) {
	gamesEndpoint := "/services/schedule_xml_2.aspx?format=json"

	if len(sports) > 0 {
		for i := 0; i < len(sports); i++ {
			gamesEndpoint += "&path=" + sports[i]
		}
	}

	if id != nil {
		gamesEndpoint += "&game_id=" + *id
	}

	if limit > 0 {
		gamesEndpoint += "&take=" + strconv.Itoa(limit)
	}

	if startDate == nil {
		chd := getChicagoTime()
		startDate = &chd
	}

	gamesEndpoint += "&starting=" + *startDate

	if endDate != nil {
		gamesEndpoint += "&ending=" + *endDate
	}

	url := host + gamesEndpoint
	bodyBytes, err := request(http.MethodGet, url, nil)

	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> GetGames: Failed to request games. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var s sidearmModel.Schedule
	es := json.Unmarshal(bodyBytes, &s)
	if es != nil {
		log.Printf("sidearm -> GetGames: Failed to unmarshal response json. Reason: %s", es.Error())
		return nil, es
	}

	games := buildGames(s)
	return games, nil
}

// GetTeamSchedule retrieves team schedule for specific year
func (p *Provider) GetTeamSchedule(sport string, year *int) (*model.Schedule, error) {
	s, err := getSportSeason(sport, year)
	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, fmt.Errorf("sidearm -> GetTeamSchedule: season was not fount")
	}

	sch, err := getSchedule(*s)
	if err != nil {
		return nil, err
	}

	games := buildGames(*sch)
	return &model.Schedule{Label: s.ScheduleYear, Games: games}, nil
}

// GetTeamRecord retrieves team record for specific year
func (p *Provider) GetTeamRecord(sport string, year *int) (*model.Record, error) {
	s, err := getSportSeason(sport, year)
	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, fmt.Errorf("sidearm -> GetTeamRecord: season was not fount")
	}

	sch, err := getSchedule(*s)
	if err != nil {
		return nil, err
	}

	r := sch.Record

	return &model.Record{OverallRecord: r.OverallRecord, ConferenceRecord: r.ConferenceRecord, Streak: r.Streak, HomeRecord: r.HomeRecord, AwayRecord: r.AwayRecord, NeutralRecord: r.NeutralRecord}, nil
}

// GetLiveGames retrieves current live games
func (p *Provider) GetLiveGames() ([]model.LiveGame, error) {
	return p.stats.LiveData(), nil
}

// GetConfig retrieves the config
func (p *Provider) GetConfig() (map[string]interface{}, error) {
	cfgBytes, err := json.Marshal(p.config)
	if err != nil {
		log.Println("sidearm -> GetConfig(): Failed to marshal config to bytes")
		return nil, err
	}

	var cfgMap map[string]interface{}
	err = json.Unmarshal(cfgBytes, &cfgMap)
	if err != nil {
		log.Println("sidearm -> GetConfig(): Failed to unmarshal config to map")
		return nil, err
	}
	return cfgMap, nil
}

// UpdateConfig updates the config
func (p *Provider) UpdateConfig(cfgBytes []byte) error {
	if cfgBytes == nil {
		msg := "new config value must not be nil"
		log.Printf("sidearm -> UpdateConfig: failed to update config. Reason: %s", msg)
		return fmt.Errorf(msg)
	}

	var cfg source.Config
	err := json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Println("sidearm -> UpdateConfig: Failed to unmarshal config bytes to struct")
		return err
	}

	p.config = cfg
	p.stats.UpdateConfig(cfg)
	return nil
}

func getSportSeason(sport string, year *int) (*sidearmModel.Season, error) {
	seasonsEndpoint := "/services/schedule_xml_2.aspx?format=json&sportseasons=true"

	seasonsEndpoint += "&path=" + sport

	if year != nil {
		seasonsEndpoint += "&year=" + strconv.Itoa(*year)
	}

	seasonsURL := host + seasonsEndpoint
	seasonsBodyBytes, err := request(http.MethodGet, seasonsURL, nil)
	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> getSportSeason: Failed to load sport seasons games. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var ss sidearmModel.SportSeasons
	err = json.Unmarshal(seasonsBodyBytes, &ss)
	if err != nil {
		log.Printf("sidearm -> getSportSeason: Failed to unmarshal response json. Reason: %s", err.Error())
		return nil, err
	}

	l := len(ss.Seasons)
	if l == 0 {
		return nil, fmt.Errorf("sidearm -> getSportSeason: there are no seasons for sport [%s]", sport)
	}

	var s *sidearmModel.Season
	if year != nil {
		// Get season for specific year
		ys := strconv.Itoa(*year)
		for i := 0; i < l; i++ {
			c := ss.Seasons[i]
			if ys == c.Year {
				s = &c
				break
			}
		}
	} else {
		// Or get the last one
		s = &ss.Seasons[l-1]
	}
	return s, nil
}

func getSchedule(s sidearmModel.Season) (*sidearmModel.Schedule, error) {
	scheduleURL := s.ScheduleURL
	if len(scheduleURL) == 0 {
		return nil, fmt.Errorf("sidearm -> getSchedule: there is no schedule url for season [%s]", s.Year)
	}

	scheduleBodyBytes, err := request(http.MethodGet, scheduleURL, nil)
	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> getSchedule: failed to load team schedule. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var sch sidearmModel.Schedule
	err = json.Unmarshal(scheduleBodyBytes, &sch)
	if err != nil {
		log.Printf("sidearm -> getSchedule: failed to unmarshal schedule response json. Reason: %s", err.Error())
		return nil, err
	}
	return &sch, nil
}

func buildGames(s sidearmModel.Schedule) []model.Game {
	var games []model.Game
	saGames := s.Games
	if (saGames != nil) && (len(saGames) > 0) {
		for i := 0; i < len(saGames); i++ {
			s := saGames[i]
			var sport model.Sport
			if s.Sport != nil {
				sport = model.Sport{Title: s.Sport.Title, ShortName: s.Sport.ShortName}
			}

			var location model.Location
			if s.Location != nil {
				location = model.Location{Location: s.Location.Location, HAN: s.Location.HAN}
			}

			var links model.Links
			if s.Links != nil {
				links = model.Links{Livestats: s.Links.Livestats, Video: s.Links.Video, Audio: s.Links.Audio, Tickets: s.Links.Tickets}
				var preGame model.GameInfo
				if s.Links.PreGame != nil {
					preGame = model.GameInfo{ID: s.Links.PreGame.ID, URL: s.Links.PreGame.URL, StoryImageURL: s.Links.PreGame.StoryImageURL, Text: s.Links.PreGame.Text}
					links.PreGame = &preGame
				}
			}

			var opponent model.Opponent
			if s.Opponent != nil {
				opponent = model.Opponent{Name: s.Opponent.Name, LogoImage: s.Opponent.LogoImage}
			}

			var results []model.Result
			if (s.Results != nil) && len(*s.Results) > 0 {
				for i := 0; i < len(*s.Results); i++ {
					r := (*s.Results)[i]
					results = append(results, model.Result{Status: r.Status, TeamScore: r.TeamScore, OpponentScore: r.OpponentScore})
				}
			}
			parkingURL := getParkingURL(s.DisplayField2)
			name := getName(s)
			games = append(games, model.Game{ID: s.ID, Name: name, Date: s.Date, DateTimeUtc: s.DateTimeUtc, EndDateTimeUtc: s.EndDateTimeUtc, EndDate: s.EndDateTime, Time: s.Time, AllDay: s.DateInfo.AllDay, Status: s.Status, Description: s.PromotionName, Sport: &sport, Location: &location, ParkingURL: parkingURL, Links: &links, Opponent: &opponent, Results: &results})
		}
	}
	return games
}

// Extracts url from "displayField2". Example:
//
// <a href="https://ev11.evenue.net/cgi-bin/ncommerce3/SEGetEventInfo?ticketCode=GS%3AILLINOIS%3AF19%3A03P%3A&linkID=illinois&shopperContext=&pc=&caller=&appCode=&groupCode=FP&cgc=&dataAccId=863&locale=en_US&siteId=ev_illinois&poolId=pac8-evcluster1&sDomain=ev11.evenue.net" target="_blank">Buy Parking</a>
func getParkingURL(displayField2 string) *string {
	if displayField2 == "" {
		return nil
	}

	hasHref := strings.HasPrefix(displayField2, "<a href=\"")
	if !hasHref {
		return nil
	}

	startQuotesIndex := strings.Index(displayField2, "\"")
	if startQuotesIndex == -1 {
		return nil
	}

	url := displayField2[(startQuotesIndex + 1):]

	endQuotesIndex := strings.Index(url, "\"")
	if endQuotesIndex == -1 {
		return nil
	}

	url = url[0:endQuotesIndex]
	return &url
}

func getChicagoTime() string {
	now := time.Now()
	tl, err := time.LoadLocation("America/Chicago")
	if err == nil {
		now = now.In(tl)
	} else {
		log.Printf("sidearm -> getChicagoTime: failed to retrieve Chicago time -> error:\n%s", err.Error())
	}
	time := now.Format("01/02/2006")
	log.Printf("sidearm -> getChicagoTime: now in Chicago:%s\tresult:%s\n", now, time)
	return time
}

func buildRosterPhotos(srcPhotos []sidearmModel.Photo) *model.Photos {
	if (srcPhotos == nil) || len(srcPhotos) <= 0 {
		return nil
	}
	var srcPhoto *sidearmModel.Photo
	for i := 0; i < len(srcPhotos); i++ {
		p := srcPhotos[i]
		// Use "headshot" photo if exists. The client use it.
		if p.Type == "headshot" {
			srcPhoto = &p
			break
		}
	}
	if srcPhoto == nil {
		srcPhoto = &srcPhotos[0]
	}
	var fsURL string
	var thURL string
	if srcPhoto.Fullsize != "" {
		fsURL = srcPhoto.Fullsize
	} else {
		thURL = srcPhoto.Roster
	}
	var photos model.Photos
	// Prepare photos with specific url for the client - either use fullsize and resized or use roster photo
	if fsURL != "" {
		photos = model.Photos{Fullsize: fsURL, Thumbnail: fmt.Sprintf("%s?width=256", fsURL)}
	} else if thURL != "" {
		photos = model.Photos{Fullsize: thURL, Thumbnail: thURL}
	}
	return &photos
}

func request(method, url string, body io.Reader) (responseBytes []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}

	// Send custom "User-Agent" header because fightingillini returns 404 Not found if "User-Agent" is not persistant
	req.Header.Set("User-Agent", "golang_sports_service")

	client := &http.Client{Transport: &http.Transport{}}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	code := resp.StatusCode
	if !((200 <= code) && (code <= 206)) {
		errMsg := string(bodyBytes)
		return nil, fmt.Errorf(errMsg)
	}

	return bodyBytes, nil
}

func (p *Provider) processCachedGames() {
	for {
		log.Println("sidearm -> processCachedGames - reload")
		p.loadCachedGames()
		p.processNextGameItems()
		// reload on each hour
		timer := time.NewTimer(time.Hour)

		<-timer.C
	}
}

func (p *Provider) loadCachedGames() {
	gamesEndpoint := fmt.Sprintf("/services/schedule_xml_2.aspx?format=json&starting=%s", getChicagoTime())

	url := host + gamesEndpoint
	bodyBytes, err := request(http.MethodGet, url, nil)

	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> loadCachedGames: Failed to request games. Reason: %s", err.Error())
		log.Print(errMsg)
		return
	}

	var schedule sidearmModel.Schedule
	err = json.Unmarshal(bodyBytes, &schedule)
	if err != nil {
		log.Printf("sidearm -> loadCachedGames: Failed to unmarshal response json. Reason: %s", err.Error())
		return
	}

	p.mu.Lock()
	p.cachedGames = schedule.Games
	p.mu.Unlock()
	log.Println("sidearm -> loadCachedGames: games loaded")
}

func (p *Provider) processNextGameItems() {

	var (
		started []*sidearmModel.LiveGameItem
		next    sidearmModel.LiveGameItem
	)

	if len(p.cachedGames) == 0 {
		return
	}

	for i := 0; i < len(p.cachedGames); i++ {
		game := p.cachedGames[i]

		gameID := game.ID
		home := getHome(game)
		opponentName := getOpponentName(game)
		livestats := getLivestats(game)

		//we need only the games with livestats
		if len(livestats) > 0 && next.Time.IsZero() {
			if game.DateTimeUtc != "" {
				if t, err := time.Parse("2006-01-02T15:04:05Z", game.DateTimeUtc); err == nil {

					if t.After(time.Now()) {
						next.GameID = gameID
						next.Time = t
						next.Sport = game.Sport.ShortName
						next.Home = home
						next.OpponentName = opponentName
					} else {
						started = append(started, &sidearmModel.LiveGameItem{GameID: gameID, Time: t, Sport: game.Sport.ShortName, Home: home, OpponentName: opponentName})
					}
				} else {
					log.Printf("sidearm -> processNextGameItems: failed to parse datetime_utc string to time. Reason: %s", err.Error())
				}
			}
		}
	}

	p.mu.Lock()
	p.nextGame = next
	p.startedGames = started
	p.mu.Unlock()
}

func (p *Provider) processLiveStats() {
	for {
		hasData := hasData(p.nextGame)

		if hasData {
			log.Println("sidearm: processLiveStats -> there is loaded games, start loading the live data...")

			items := p.startedGames
			items = append(items, &p.nextGame) //add next
			err := p.stats.ProcessLiveData(items)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println("sidearm: processLiveStats -> there is no loaded games yet, so we need to wait for loading the live data...")
		}

		//Calculate the interval duration
		duration := p.calculateIntervalDuration()
		//Trigger timer for next processing
		timer := time.NewTimer(duration)

		<-timer.C
	}
}

func (p *Provider) calculateIntervalDuration() time.Duration {
	isDuringLive := p.stats.IsDuringLiveGame()
	if isDuringLive {
		log.Printf("sidearm: calculateIntervalDuration -> During live")
		return time.Second * 3 //process every 3 seconds
	}
	if p.isPreGame() {
		log.Printf("sidearm: calculateIntervalDuration -> Pre game")
		return time.Second * 5 //process every 5 seconds
	}
	duration := time.Hour
	nextItem := p.nextGame
	if hasData(nextItem) {
		log.Printf("sidearm: calculateIntervalDuration -> Next - gameid:%s\tsport:%s\thome:%t\tstart:%s\n", nextItem.GameID, nextItem.Sport, nextItem.Home, nextItem.Time)
		now := time.Now()
		nextDelta := nextItem.Time.Add(time.Duration(-5) * time.Minute) //5 minutes before the game
		startDuration := nextDelta.Sub(now)
		if startDuration > 0 && startDuration < duration {
			duration = startDuration
		}
	} else {
		log.Println("sidearm: calculateIntervalDuration -> there is no next yet")
		duration = time.Second * 5 //process every 5 seconds, we need to know the next start
	}
	log.Printf("sidearm: calculateIntervalDuration -> Next processing after:%s\n", duration)
	return duration
}

func (p *Provider) isPreGame() bool {
	//It is pre game if it is 5 minutes before the start time of the game
	next := p.nextGame.Time
	if next.IsZero() {
		return false
	}

	nowInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	nextInMilliSeconds := next.UnixNano() / int64(time.Millisecond)
	preGameStartInMilliSeconds := nextInMilliSeconds - int64(5*60*1000) // - 5 minutes

	if nowInMilliSeconds >= preGameStartInMilliSeconds {
		return true
	}
	return false
}

func (p *Provider) loadNews(id *string, sports []string, limit int) ([]model.News, error) {
	newsEndpoint := "/services/stories_xml.aspx?format=json"
	if id != nil {
		newsEndpoint += "&story_id=" + *id
	}

	if len(sports) > 0 {
		for i := 0; i < len(sports); i++ {
			newsEndpoint += "&path=" + sports[i]
		}
	}

	if limit > 0 {
		newsEndpoint += "&numrec=" + strconv.Itoa(limit)
	}

	url := host + newsEndpoint
	bodyBytes, err := request(http.MethodGet, url, nil)

	if err != nil {
		errMsg := fmt.Sprintf("sidearm -> loadNews: Failed to request news. Reason: %s", err.Error())
		log.Print(errMsg)
		return nil, err
	}

	var s sidearmModel.Stories
	em := json.Unmarshal(bodyBytes, &s)
	if em != nil {
		log.Printf("sidearm -> loadNews: Failed to unmarshal response json. Reason: %s", em.Error())
		return nil, em
	}

	stories := s.Stories
	if stories == nil {
		return nil, nil
	}

	var news []model.News
	storiesLength := len(stories)

	if storiesLength > 0 {
		for i := 0; i < storiesLength; i++ {
			s := stories[i]
			var e = &s.Enclosure
			var sport = &s.Sport
			var imageURL string
			if e != nil {
				imageURL = e.URL
			}
			news = append(news, model.News{ID: s.ID, Title: s.Title, Sport: sport.PrimaryGlobalShortName, Link: s.Link, Category: s.Category, Description: s.Description, FullText: s.FullText, FullTextRaw: s.FullTextRaw, ImageURL: imageURL, PubDateUtc: s.PubDateUtc})
		}
	}

	return news, nil
}

func (p *Provider) processCachedNews() {
	for {
		p.checkUpdatedNews()
		timer := time.NewTimer(time.Minute)

		<-timer.C
	}
}

func (p *Provider) loadCachedNews() {
	news, err := p.loadNews(nil, nil, 0)
	if err != nil {
		log.Printf("sidearm -> loadCachedNews: Failed to load cached news. Reason: %s", err.Error())
		return
	}

	p.mu.Lock()
	p.cachedNews = news
	p.mu.Unlock()
	log.Println("sidearm -> loadCachedNews: success")
}

func (p *Provider) checkUpdatedNews() {
	curNews, err := p.loadNews(nil, nil, 0)
	if err != nil {
		log.Printf("sidearm -> checkUpdatedNews: Failed to load cached news. Reason: %s", err.Error())
		return
	}

	if len(curNews) == 0 {
		log.Println("sidearm -> checkUpdatedNews: No current news.")
		return
	}

	if len(p.cachedNews) == 0 {
		log.Println("sidearm -> checkUpdatedNews: No cached news.")
		return
	}

	for i := 0; i < len(curNews); i++ {
		item := curNews[i]
		exists := p.isNewsCached(item)
		if !exists {
			log.Printf("sidearm -> checkUpdatedNews: Found new item: %s", item.ID)
			p.sendNewsNotification(item)
		}
	}

	p.mu.Lock()
	p.cachedNews = curNews
	p.mu.Unlock()
}

func (p *Provider) sendNewsNotification(news model.News) {
	configGameMessages := p.config.NotificationConfig.Messages
	category := news.Category
	var msgTitle string
	if len(category) > 0 {
		titleFmt := configGameMessages["news_updates_sport_title_format"]
		msgTitle = fmt.Sprintf(titleFmt, category)
	} else {
		msgTitle = configGameMessages["news_updates_default_title"]
	}

	msgBodyFormat := configGameMessages["news_updates_body_content_format"]
	msgBody := fmt.Sprintf(msgBodyFormat, news.Title)

	// topic is "athletics.{sport_short_name}.notification.news"
	topic := fmt.Sprintf("athletics.%s.notification.news", news.Sport)

	data := make((map[string]string))
	data["type"] = "athletics_news_detail"
	data["sport"] = news.Sport
	data["news_id"] = news.ID
	data["click_action"] = "FLUTTER_NOTIFICATION_CLICK"
	err := p.notifications.SendNotificationMsg(topic, msgTitle, msgBody, data)
	if err != nil {
		log.Printf("sidearm -> sendNewsNotification: error sending notification topic:%s title:%s body:%s data:%s %s", topic, msgTitle, msgBody, data, err.Error())
	} else {
		log.Printf("sidearm -> sendNewsNotification: success sending notification topic:%s title:%s body:%s data:%s", topic, msgTitle, msgBody, data)
	}
}

func (p *Provider) isNewsCached(item model.News) bool {
	for _, b := range p.cachedNews {
		if b.ID == item.ID {
			return true
		}
	}
	return false
}

func hasData(item sidearmModel.LiveGameItem) bool {
	return len(item.GameID) > 0 && len(item.Sport) > 0 && !item.Time.IsZero()
}

func getHome(scheduleItem sidearmModel.Game) bool {
	var result bool
	if scheduleItem.Location != nil {
		han := scheduleItem.Location.HAN
		result = "H" == han
	}
	return result
}

func getOpponentName(scheduleItem sidearmModel.Game) string {
	var result string
	if scheduleItem.Opponent != nil {
		result = scheduleItem.Opponent.Name
	}
	return result
}

func getLivestats(scheduleItem sidearmModel.Game) string {
	var result string
	if scheduleItem.Links != nil {
		result = scheduleItem.Links.Livestats
	}
	return result
}

func getName(game sidearmModel.Game) string {
	var name string
	var opponentName = getOpponentName(game)
	if getHome(game) {
		name = opponentName + " at " + illinoisTeamName
	} else {
		name = illinoisTeamName + " at " + opponentName
	}
	if game.Status == "C" {
		name = name + " (Cancelled)"
	}
	return name
}
