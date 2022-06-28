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
