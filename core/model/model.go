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

package model

// News structure
type News struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Category    string `json:"category"`
	Sport       string `json:"sport"`
	Description string `json:"description"`
	FullText    string `json:"fulltext"`
	FullTextRaw string `json:"fulltext_raw"`
	ImageURL    string `json:"image_url"`
	PubDateUtc  string `json:"pub_date_utc"`
}

// Coach structure
type Coach struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	FirstName string  `json:"first_name,omitempty"`
	LastName  string  `json:"last_name,omitempty"`
	Email     string  `json:"email,omitempty"`
	Phone     string  `json:"phone,omitempty"`
	Title     string  `json:"title,omitempty"`
	Bio       string  `json:"bio,omitempty"`
	Photos    *Photos `json:"photos,omitempty"`
}

// Player structure
type Player struct {
	ID         string  `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	FirstName  string  `json:"first_name,omitempty"`
	LastName   string  `json:"last_name,omitempty"`
	Uni        string  `json:"uni,omitempty"`
	PosShort   string  `json:"pos_short,omitempty"`
	Height     string  `json:"height,omitempty"`
	Wight      string  `json:"weight,omitempty"`
	Gender     string  `json:"gender,omitempty"`
	YearLong   string  `json:"year_long,omitempty"`
	HomeTown   string  `json:"hometown,omitempty"`
	HighSchool string  `json:"highschool,omitempty"`
	Captain    bool    `json:"captain,omitempty"`
	Bio        string  `json:"bio,omitempty"`
	Photos     *Photos `json:"photos,omitempty"`
}

// Photos structure
type Photos struct {
	Fullsize  string `json:"fullsize,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// SportSocial structure
type SportSocial struct {
	OrgID          string `json:"org_id"`
	SportShortName string `json:"shortname,omitempty"`
	TwitterName    string `json:"sport_twitter_name,omitempty"`
	InstagramName  string `json:"sport_instagram_name,omitempty"`
	FacebookPage   string `json:"sport_facebook_page,omitempty"`
}

// Game structure
type Game struct {
	ID             string    `json:"id,omitempty"`
	Date           string    `json:"date,omitempty"`
	DateTimeUtc    string    `json:"datetime_utc,omitempty"`
	EndDate        string    `json:"end_date,omitempty"`
	EndDateTimeUtc string    `json:"end_datetime_utc,omitempty"`
	Time           string    `json:"time,omitempty"`
	AllDay         bool      `json:"all_day,omitempty"`
	Status         string    `json:"status,omitempty"`
	Description    string    `json:"description,omitempty"`
	Sport          *Sport    `json:"sport,omitempty"`
	Location       *Location `json:"location,omitempty"`
	ParkingURL     *string   `json:"parking_url,omitempty"`
	Links          *Links    `json:"links,omitempty"`
	Opponent       *Opponent `json:"opponent,omitempty"`
	Results        *[]Result `json:"results,omitempty"`
}

// Sport structure
type Sport struct {
	Title     string `json:"title,omitempty"`
	ShortName string `json:"shortname,omitempty"`
}

// Location structure
type Location struct {
	Location string `json:"location,omitempty"`
	HAN      string `json:"HAN,omitempty"`
}

// Links structure
type Links struct {
	Livestats string    `json:"livestats,omitempty"`
	Video     string    `json:"video,omitempty"`
	Audio     string    `json:"audio,omitempty"`
	Tickets   string    `json:"tickets,omitempty"`
	PreGame   *GameInfo `json:"pregame,omitempty"`
}

// GameInfo structure
type GameInfo struct {
	ID            string `json:"id,omitempty"`
	URL           string `json:"url,omitempty"`
	StoryImageURL string `json:"story_image_url,omitempty"`
	Text          string `json:"text,omitempty"`
}

// Opponent structure
type Opponent struct {
	Name      string `json:"name,omitempty"`
	LogoImage string `json:"logo_image,omitempty"`
}

// Result structure
type Result struct {
	Status        string `json:"status,omitempty"`
	TeamScore     string `json:"team_score,omitempty"`
	OpponentScore string `json:"opponent_score,omitempty"`
}

// Schedule structure
type Schedule struct {
	Label string `json:"label,omitempty"`
	Games []Game `json:"games,omitempty"`
}

// Record structure
type Record struct {
	OverallRecord    string `json:"overall_record_unformatted,omitempty"`
	ConferenceRecord string `json:"conference_record,omitempty"`
	Streak           string `json:"streak,omitempty"`
	HomeRecord       string `json:"home_record,omitempty"`
	AwayRecord       string `json:"away_record,omitempty"`
	NeutralRecord    string `json:"neutral_record,omitempty"`
}

// LiveGame interface
type LiveGame interface {
	GetType() string
	GetGameID() int
	GetPath() string
	GetHasStarted() bool
	GetIsComplete() bool
	GetClockSeconds() int
	GetPeriod() int
	GetHomeScore() int
	GetVisitingScore() int

	GetCustomData() string //every sport can add custom data, football add possession and last play for example

	Encode() map[string]string
}
