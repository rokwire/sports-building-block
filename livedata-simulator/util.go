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

func fillInitialData() []game {
	home1 := team{
		ID:    "Illini",
		Name:  "Illinois",
		Score: 0,
	}

	visiting1 := team{
		ID:    "Flyers",
		Name:  "Dayton",
		Score: 0,
	}

	game1 := game{
		GameID:            16799,
		Path:              "wsoc",
		FullPath:          "http://static.sidearmstats.com/schools/illinois/wsoc/",
		Link:              "https://sidearmstats.com/illinois/wsoc/",
		SportTitle:        "Soccer",
		Opponent:          "Dayton",
		Location:          "Champaign, Ill.",
		Time:              "1 PM",
		HasStarted:        false,
		IsComplete:        false,
		ClockSeconds:      0,
		PeriodsRegulation: 2,
		Period:            1,
		HomeTeam:          home1,
		VisitingTeam:      visiting1,
	}

	Games := make([]game, 0, 1)
	Games = append(Games, game1)
	return Games
}
