package source

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sport/core/model"
	sidearmModel "sport/driven/provider/sidearm/model"
	"strconv"
)

var statsURL = "https://fightingillini.com/services/livestats.ashx"

type sidearmGames struct {
	Games []sidearmGame
}

type sidearmTeam struct {
	ID    string `json:"Id"`
	Name  string
	Score int
}

type sidearmGame struct {
	GameID            int `json:"GameId"`
	Path              string
	FullPath          string
	Link              string
	SportTitle        string
	Opponent          string
	Location          string
	Time              string
	HasStarted        bool
	IsComplete        bool
	ClockSeconds      int
	PeriodsRegulation int
	Period            int
	HomeTeam          sidearmTeam
	VisitingTeam      sidearmTeam
}

func (game *sidearmGame) GetType() string {
	return game.Path
}

func (game *sidearmGame) GetGameID() int {
	return game.GameID
}

func (game *sidearmGame) GetPath() string {
	return game.Path
}

func (game *sidearmGame) GetHasStarted() bool {
	return game.HasStarted
}

func (game *sidearmGame) GetIsComplete() bool {
	return game.IsComplete
}

func (game *sidearmGame) GetClockSeconds() int {
	return game.ClockSeconds
}

func (game *sidearmGame) GetPeriod() int {
	return game.Period
}

func (game *sidearmGame) GetHomeScore() int {
	return game.HomeTeam.Score
}

func (game *sidearmGame) GetVisitingScore() int {
	return game.VisitingTeam.Score
}

func (game *sidearmGame) GetCustomData() string {
	return ""
}

func (game *sidearmGame) Encode() map[string]string {

	data := make((map[string]string))
	data["Type"] = game.GetType()
	data["GameId"] = strconv.Itoa(game.GameID)
	data["Path"] = game.Path
	data["HasStarted"] = strconv.FormatBool(game.HasStarted)
	data["IsComplete"] = strconv.FormatBool(game.IsComplete)
	data["ClockSeconds"] = strconv.Itoa(game.ClockSeconds)
	data["Period"] = strconv.Itoa(game.Period)
	data["HomeScore"] = strconv.Itoa(game.HomeTeam.Score)
	data["VisitingScore"] = strconv.Itoa(game.VisitingTeam.Score)
	data["Custom"] = game.GetCustomData()

	return data
}

type sidearmSource struct {
	config Config
}

func newSidearmSource(config Config) sidearmSource {
	var sidearmSource sidearmSource
	sidearmSource.config = config
	return sidearmSource
}

func (sidearmSource *sidearmSource) updateConfig(config Config) {
	log.Println("sidearmsports: UpdateConfig -> config updated in sidearm source")
	sidearmSource.config = config
}

func (sidearmSource *sidearmSource) load(item *sidearmModel.LiveGameItem) (model.LiveGame, error) {
	var (
		b      []byte
		r      *http.Request
		resp   *http.Response
		err    error
		games  sidearmGames
		result *sidearmGame
	)

	client := &http.Client{Transport: &http.Transport{}}

	r, err = http.NewRequest(http.MethodGet, statsURL, nil)

	if err == nil {
		//Send custom "User-Agent" header because fightingillini returns 404 Not found if "User-Agent" is not persistant otherwise
		r.Header.Set("User-Agent", "golang_sports_service")
		resp, err = client.Do(r)
	}

	if err == nil {
		b, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	if err != nil {
		b = nil
	}

	if b != nil {
		err = json.Unmarshal(b, &games)
	}

	if err != nil {
		log.Printf("sidearmsports: loadFromSideArm -> Error loading live data:%s\n", err.Error())
		return nil, err
	}
	if len(games.Games) == 0 {
		return nil, errors.New("sidearmsports: loadFromSideArm -> No games")
	}

	for _, element := range games.Games {
		if strconv.Itoa(element.GameID) == item.GameID {
			result = &element
		}
	}
	if result == nil {
		return nil, errors.New("sidearmsports: loadFromSideArm -> there is games but not the presented one")
	}

	return result, nil
}
