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
