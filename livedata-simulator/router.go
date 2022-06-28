package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type routerHandler struct {
	application *application
}

func newRouter(app *application) routerHandler {
	return routerHandler{
		application: app,
	}
}

func (rH routerHandler) setRoutes() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/start", rH.start).Methods("GET")
	router.HandleFunc("/stop", rH.stop).Methods("GET")
	router.HandleFunc("/data", rH.data).Methods("GET")
	http.ListenAndServe(":8585", router)
}

func (rH routerHandler) start(w http.ResponseWriter, r *http.Request) {
	rH.application.started = true

	go start(rH.application)
	w.WriteHeader(http.StatusOK)
}

func start(app *application) {
	//1. fill the initial data
	app.Games = fillInitialData()

	//start simulation timer
	starSimulationTimer(app)
}

func starSimulationTimer(app *application) {
	for app.started {
		updateInterval := rand.Intn(4) + 3          //random 3 <= n <=6
		app.simulationTimeSeconds += updateInterval // add the update interval to the simulation time

		//first two minutes in pregame
		if 0 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 120 {
			applyPregame(app, updateInterval)
		} else if 120 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 240 {
			//the next two minutes is the 1st half
			applyFirstHalf(app, updateInterval)
		} else if 240 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 360 {
			//the next one minute is the half time
			applyHalfTime(app, updateInterval)
		} else if 360 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 600 {
			//the next two minutes is the 2nd half
			applySecondHalf(app, updateInterval)
		} else if 600 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 840 {
			//the next two minutes is after game
			applyAfterGame(app, updateInterval)
		}

		interval := time.Second * time.Duration(updateInterval)
		timer := time.NewTimer(interval)
		select {
		case <-app.exit:
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func applyPregame(app *application, updateInterval int) {
	log.Printf("Pregame:%d", app.simulationTimeSeconds)

	//do nothing
	printGame(app)
}

func applyFirstHalf(app *application, updateInterval int) {
	log.Printf("1st Half:%d", app.simulationTimeSeconds)

	game1 := app.Games[0]

	//mark as started
	if !game1.HasStarted {
		game1.HasStarted = true
	}

	//set the period
	if game1.Period != 1 {
		game1.Period = 1
	}

	// add the update interval to the clock
	game1.ClockSeconds = game1.ClockSeconds + updateInterval

	if 140 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 146 {
		// add goal to the home team
		homeTeam := game1.HomeTeam
		homeTeam.Score = 1
		game1.HomeTeam = homeTeam
	}

	app.Games[0] = game1

	printGame(app)
}

func applyHalfTime(app *application, updateInterval int) {
	log.Printf("Half Time:%d", app.simulationTimeSeconds)

	//do nothing

	printGame(app)
}

func applySecondHalf(app *application, updateInterval int) {
	log.Printf("2nd Half:%d", app.simulationTimeSeconds)

	game1 := app.Games[0]

	//set the period
	if game1.Period != 2 {
		game1.Period = 2
	}

	// add the update interval to the clock
	game1.ClockSeconds = game1.ClockSeconds + updateInterval

	if 400 <= app.simulationTimeSeconds && app.simulationTimeSeconds < 406 {
		// add goal to the visiting team
		visitingTeam := game1.VisitingTeam
		visitingTeam.Score = 1
		game1.VisitingTeam = visitingTeam
	}

	app.Games[0] = game1

	printGame(app)
}

func applyAfterGame(app *application, updateInterval int) {
	log.Printf("After game:%d", app.simulationTimeSeconds)

	game1 := app.Games[0]

	//mark as completed
	if !game1.IsComplete {
		game1.IsComplete = true
	}

	app.Games[0] = game1

	printGame(app)
}

func printGame(app *application) {
	if app.Games == nil {
		log.Printf("nil\n")
	} else {
		game := app.Games[0]
		log.Printf("[Clock:%d\tPeriod:%d\tStarted:%v\tCompleted:%v\tHome Score:%b\tVisiting Score:%b]\n",
			game.ClockSeconds, game.Period, game.HasStarted, game.IsComplete, game.HomeTeam.Score, game.VisitingTeam.Score)
	}
}

func (rH routerHandler) stop(w http.ResponseWriter, r *http.Request) {
	app := rH.application
	app.started = false
	app.Games = nil
	app.simulationTimeSeconds = 0

	close(rH.application.exit)
	w.WriteHeader(http.StatusOK)
}

func (rH routerHandler) data(w http.ResponseWriter, r *http.Request) {
	games := rH.application.Games
	var data []byte
	if games == nil {
		data = []byte("{ \"Games\": []}")
	} else {
		//simulate the inconsistent behaviour of the API
		//during the game just return empty data
		inconsistentFlag := rand.Intn(3)
		if inconsistentFlag == 0 {
			data = []byte("{ \"Games\": []}")
		} else {
			gamesDataBytes, _ := json.Marshal(rH.application.Games)
			gamesDataStr := string(gamesDataBytes)
			data = []byte("{ \"Games\": " + gamesDataStr + "}")
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
