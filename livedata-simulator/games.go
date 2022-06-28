package main

type game struct {
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
	HomeTeam          team
	VisitingTeam      team
}

type team struct {
	ID    string `json:"Id"`
	Name  string
	Score int
}

type games struct {
	Games []game
}
