package core

import (
	"sport/core/model"
)

// Storage interface has to be implemented by all Storage adapters
type Storage interface {
	GetSports()
}

// Provider interface has to be implemented by all sports providers
type Provider interface {
	GetNews(id *string, sport []string, limit int) ([]model.News, error)
	GetCoaches(sport string) ([]model.Coach, error)
	GetPlayers(sport string) ([]model.Player, error)
	GetSocialNetworks() ([]model.SportSocial, error)
	GetGames(sports []string, id *string, startDate *string, endDate *string, limit int) ([]model.Game, error)
	GetTeamSchedule(sport string, year *int) (*model.Schedule, error)
	GetTeamRecord(sport string, year *int) (*model.Record, error)
	GetLiveGames() ([]model.LiveGame, error)
	GetConfig() (map[string]interface{}, error)
	UpdateConfig(data []byte) error
}
