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

package livestats

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sport/core/model"
	"sport/driven/notifications"
	"sport/driven/provider/sidearm/livestats/source"
	sidearmModel "sport/driven/provider/sidearm/model"
	"strconv"
)

// LiveStats service
type LiveStats interface {
	UpdateConfig(config source.Config)
	ProcessLiveData(items []*sidearmModel.LiveGameItem) error
	IsDuringLiveGame() bool
	LiveData() []model.LiveGame
}

type livestats struct {
	config        source.Config
	notifications notifications.Notifications
	games         sidearmModel.GameItems
	lsSource      source.Source
	teamName      string
}

// New create live stats checker
func New(notifications notifications.Notifications, config source.Config, ftpHost string, ftpUser string, ftpPassword string, teamName string) LiveStats {
	lsSource := source.New(config, ftpHost, ftpUser, ftpPassword)
	return &livestats{config: config, notifications: notifications, lsSource: lsSource, teamName: teamName}
}

func (stats *livestats) UpdateConfig(config source.Config) {
	log.Println("LiveStats: UpdateConfig -> config updated in livestats")
	stats.config = config
	stats.lsSource.UpdateConfig(config)
}

func (stats *livestats) ProcessLiveData(items []*sidearmModel.LiveGameItem) error {
	if len(items) <= 0 {
		return errors.New("LiveStats: ProcessLiveData -> cannot load the live data for nil items")
	}

	for _, item := range items {
		stats.processLiveDataForItem(item)
	}
	return nil
}

func (stats *livestats) processLiveDataForItem(item *sidearmModel.LiveGameItem) {
	if item == nil {
		log.Println("LiveStats: processLiveDataForItem -> cannot process data for nil item")
	}
	//1. load the live data
	loadedGameItem, err := stats.lsSource.LoadData(item)
	if err != nil {
		sport := item.Sport
		gameID := item.GameID
		home := item.Home
		sources := stats.config.GetLivestatsSource(sport, home)
		log.Printf("LiveStats: processLiveDataForItem -> cannot load live data for item:%s %s %s error %s\n", sources, sport, gameID, err.Error())
		return
	}

	gameID := loadedGameItem.GetGameID()
	log.Printf("LiveStats: processLiveDataForItem -> the game item was loaded: %d\n", gameID)

	//2. check if we need game changed notification
	needsGameChangedNotification := stats.needsGameChangedNotification(loadedGameItem)

	//3. check if we need game state changed notification
	needsGameStateChangedNotification, gameStarted := stats.needsGameStateChangedNotification(loadedGameItem)

	//4. update it to the list
	foundedGameItem, foundIndex := stats.findGame(loadedGameItem)
	if foundedGameItem == nil {
		// add item
		stats.games.Games = append(stats.games.Games, loadedGameItem)
		log.Printf("LiveStats: processLiveDataForItem -> the game item was added - %d\n", gameID)
	} else {
		// update item
		stats.games.Games[foundIndex] = loadedGameItem
		log.Printf("LiveStats: processLiveDataForItem -> the game item was updated - %d\n", gameID)
	}

	//5. send game changed notification if needed
	if needsGameChangedNotification {
		log.Printf("sidearm: processLiveDataForItem -> needs game changed notification - %d\n", gameID)
		stats.notifyGameChanged(loadedGameItem)
	} else {
		log.Printf("sidearm: processLiveDataForItem -> do not need game changed notification - %d\n", gameID)
	}

	//6. send game state changed notification if needed
	if needsGameStateChangedNotification {
		log.Printf("LiveStats: processLiveDataForItem -> needs user notification - %d\n", gameID)
		stats.notifyGameStateChanged(loadedGameItem, item, gameStarted)
	} else {
		log.Printf("LiveStats: processLiveDataForItem -> do not need user notification - %d\n", gameID)
	}
}

func (stats *livestats) LiveData() []model.LiveGame {
	return stats.games.Games
}

func (stats *livestats) IsDuringLiveGame() bool {
	if stats.games.Games == nil || len(stats.games.Games) == 0 {
		//no games
		return false
	}
	for _, game := range stats.games.Games {
		//check for started but not completed yet
		if game.GetHasStarted() && !game.GetIsComplete() {
			return true
		}
	}
	return false
}

func (stats *livestats) needsGameChangedNotification(newGame model.LiveGame) bool {
	foundedGame, _ := stats.findGame(newGame)
	if foundedGame == nil {
		//it needs game changed notification if this is the first data which comes
		return true
	}
	equals := reflect.DeepEqual(foundedGame, newGame)
	if !equals {
		//it needs game changed notification if they are not equals
		return true
	}
	return false
}

func (stats *livestats) needsGameStateChangedNotification(newGame model.LiveGame) (needsNotification bool, gameStarted bool) {
	foundedGame, _ := stats.findGame(newGame)
	newGameStarted := (foundedGame == nil) && (newGame != nil) && (newGame.GetHasStarted() == true)
	if newGameStarted {
		//It needs game state notification if we receive new game with hasStarted = true
		return true, true
	} else if foundedGame == nil {
		//It does not need game state notification if new game is not started
		return false, false
	}
	hasJustStarted := (((foundedGame.GetHasStarted() == false) && (foundedGame.GetIsComplete() == false)) && ((newGame.GetHasStarted() == true) && (newGame.GetIsComplete() == false)))
	if hasJustStarted {
		//it needs game state changed notification because the game has just started
		return true, true
	}
	hasJustFinished := ((foundedGame.GetHasStarted() == true) && (newGame.GetHasStarted() == true) && (foundedGame.GetIsComplete() == false) && (newGame.GetIsComplete() == true))
	if hasJustFinished {
		//it needs game state notification because the game has just finished
		return true, false
	}
	return false, false
}

func (stats *livestats) findGame(newGame model.LiveGame) (model.LiveGame, int) {
	if stats.games.Games == nil || len(stats.games.Games) == 0 {
		return nil, -1
	}
	for index, game := range stats.games.Games {
		if game.GetGameID() == newGame.GetGameID() {
			return game, index
		}
	}
	return nil, -1
}

func (stats *livestats) notifyGameChanged(game model.LiveGame) {
	path := game.GetPath()
	data := game.Encode()
	err := stats.notifications.SendDataMsg(path, data)
	if err != nil {
		log.Printf("LiveStats: notifyGameChanged -> error sending notification topic:%s data:%s %s", path, data, err.Error())
	} else {
		log.Printf("LiveStats: notifyGameChanged -> success sending notification topic:%s data:%s", path, data)
	}
}

// build notification using configuration and send
func (stats *livestats) notifyGameStateChanged(game model.LiveGame, item *sidearmModel.LiveGameItem, started bool) {
	homeTeam, visitingTeam := stats.getTeamNames(*item)
	configGameMessages := stats.config.NotificationConfig.Messages
	msgTitleFormat := configGameMessages["game_title_format"]
	if len(msgTitleFormat) == 0 {
		msgTitleFormat = ""
	}
	// build notification title
	title := fmt.Sprintf(msgTitleFormat, homeTeam, visitingTeam)

	// build body content
	var body string
	var gameState string
	if started {
		gameState = "start"
		// Game started content
		body = configGameMessages["game_started_msg"]
	} else {
		gameState = "end"
		// Game ended content
		gameEndedMsg := configGameMessages["game_ended_msg"]
		if len(gameEndedMsg) == 0 {
			gameEndedMsg = ""
		}

		gameEndedFormat := configGameMessages["game_ended_score_format"]
		if len(gameEndedFormat) == 0 {
			gameEndedFormat = ""
		}
		gameEndedScore := fmt.Sprintf(gameEndedFormat, homeTeam, game.GetHomeScore(), visitingTeam, game.GetVisitingScore())
		body = fmt.Sprintf("%s %s", gameEndedMsg, gameEndedScore)
	}
	// topic is "athletics.{path}.notification.{start|end}"
	topic := fmt.Sprintf("athletics.%s.notification.%s", game.GetPath(), gameState)
	data := make((map[string]string))
	data["GameId"] = strconv.Itoa(game.GetGameID())
	data["Path"] = game.GetPath()
	data["click_action"] = "FLUTTER_NOTIFICATION_CLICK"
	err := stats.notifications.SendNotificationMsg(topic, title, body, data)
	if err != nil {
		log.Printf("LiveStats: notifyGameStateChanged -> error sending notification topic:%s title:%s body:%s data:%s %s", topic, title, body, data, err.Error())
	} else {
		log.Printf("LiveStats: notifyGameStateChanged -> success sending notification topic:%s title:%s body:%s data:%s", topic, title, body, data)
	}
}

func (stats *livestats) getTeamNames(item sidearmModel.LiveGameItem) (string, string) {
	if item.Home {
		return stats.teamName, item.OpponentName
	}
	return item.OpponentName, stats.teamName
}
