package model

import (
	"sport/core/model"
	"time"
)

// Stories structure
type Stories struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Stories     []Story `json:"stories"`
}

// Story structure
type Story struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Category    string     `json:"category"`
	Sport       StorySport `json:"sport"`
	Description string     `json:"description"`
	FullText    string     `json:"fulltext"`
	FullTextRaw string     `json:"fulltext_raw"`
	Enclosure   Enclosure  `json:"enclosure"`
	PubDateUtc  string     `json:"pubDateUTC"`
}

// Enclosure structure
type Enclosure struct {
	URL string `json:"url"`
}

// StorySport structure
type StorySport struct {
	PrimaryGlobalShortName string `json:"primary_global_sport_shortname"`
}

// Rosters structure
type Rosters struct {
	Rosters []Roster `json:"roster"`
}

// Roster structure
type Roster struct {
	Name       string     `json:"name"`
	FirstName  string     `json:"firstname"`
	LastName   string     `json:"lastname"`
	StaffID    int        `json:"staff_id"`
	RcID       int        `json:"rc_id"`
	RpID       string     `json:"rp_id"`
	PlayerID   string     `json:"player_id"`
	StaffInfo  StaffInfo  `json:"staffinfo"`
	PlayerInfo PlayerInfo `json:"playerinfo"`
	Bio        string     `json:"bio"`
	Photos     []Photo    `json:"photos"`
}

// StaffInfo structure
type StaffInfo struct {
	Title   string `json:"title"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	BioLink string `json:"biolink"`
}

// PlayerInfo structure
type PlayerInfo struct {
	Uni            string   `json:"uni"`
	Uni2           string   `json:"uni_2"`
	PosLong        string   `json:"pos_long"`
	PosShort       string   `json:"pos_short"`
	PosShortList   []string `json:"pos_short_list"`
	Height         string   `json:"height"`
	HeightFeet     string   `json:"height_feet"`
	HeightInches   string   `json:"height-inches"`
	Weight         string   `json:"weight"`
	Gender         string   `json:"gender"`
	Year           string   `json:"year"`
	YearLong       string   `json:"year_long"`
	HomeTown       string   `json:"hometown"`
	HighSchool     string   `json:"highschool"`
	PrevSchool     string   `json:"previous_school"`
	Major          string   `json:"major"`
	FbUsrName      string   `json:"facebook_username"`
	InstUsrName    string   `json:"instagram_username"`
	TwitUsrName    string   `json:"twitter_username"`
	SnapUsrName    string   `json:"snapchat_username"`
	TikTokUsrName  string   `json:"tiktok_username"`
	CameoUsrName   string   `json:"cameo_username"`
	YouTubeUsrName string   `json:"youtube_username"`
	TwitchUsrName  string   `json:"twitch_username"`
	ShopURL        string   `json:"shop_url"`
	BioLink        string   `json:"biolink"`
	Captain        string   `json:"captain"`
}

// Photo structure
type Photo struct {
	Type      string `json:"type"`
	Fullsize  string `json:"fullsize"`
	Roster    string `json:"roster"`
	Thumbnail string `json:"thumbnail"`
}

// SportsSocNet structure
type SportsSocNet struct {
	SportsSocial []SocNets `json:"sports"`
}

// SocNets structure
type SocNets struct {
	ID             int    `json:"id"`
	ShortName      string `json:"shortname"`
	Abbrev         string `json:"abbrev"`
	Name           string `json:"name"`
	DisplayName    string `json:"short_display_name"`
	SportShortName string `json:"global_sport_shortname"`
	ConfTitle      string `json:"conference_title"`
	ConfWebsite    string `json:"conference_website"`
	ConfID         string `json:"global_conference_id"`
	ConfDevision   string `json:"global_conference_division"`
	TwitterName    string `json:"sport_twitter_name"`
	InstagramName  string `json:"sport_instagram_name"`
	FacebookPage   string `json:"sport_facebook_page"`
	FacebookID     string `json:"sport_facebook_id"`
}

// Schedule structure
type Schedule struct {
	Games  []Game `json:"schedule"`
	Record Record `json:"record"`
}

// Game structure
type Game struct {
	ID              string            `json:"id"`
	Date            string            `json:"date"`
	FormattedDate   string            `json:"formatted_date"`
	DateTimeUtc     string            `json:"datetime_utc"`
	EndDateTimeUtc  string            `json:"end_datetime_utc"`
	EndDateTime     string            `json:"end_datetime"`
	DateInfo        *DateInfo         `json:"date_info"`
	Time            string            `json:"time"`
	Type            string            `json:"type"`
	Status          string            `json:"status"`
	NoPlayText      string            `json:"noplay_text"`
	DoubleHeader    string            `json:"doubleheader"`
	PromotionName   string            `json:"game_promotion_name"`
	PromotionLink   string            `json:"game_promotion_link"`
	Sport           *Sport            `json:"sport"`
	Location        *Location         `json:"location"`
	TV              string            `json:"tv"`
	TvImage         string            `json:"tv_image"`
	Radio           string            `json:"radio"`
	DisplayField1   string            `json:"custom_display_field_1"`
	DisplayField2   string            `json:"custom_display_field_2"`
	DisplayField3   string            `json:"custom_display_field_3"`
	TMEventID       string            `json:"ticketmaster_event_id"`
	Links           *Links            `json:"links"`
	Opponent        *Opponent         `json:"opponent"`
	Sponsor         string            `json:"sponsor"`
	Results         *[]Result         `json:"results"`
	AllAccessVideos *[]AllAccessVideo `json:"allaccess_videos"`
}

// DateInfo structure
type DateInfo struct {
	Tbd           bool   `json:"tbd"`
	AllDay        bool   `json:"all_day"`
	StartDateTime string `json:"start_datetime"`
	StartDate     string `json:"start_date"`
}

// Sport structure
type Sport struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ShortName   string `json:"shortname"`
	DisplayName string `json:"sport_short_display"`
	Abbrev      string `json:"abbrev"`
	GlobalID    int    `json:"global_sport_id"`
}

// Location structure
type Location struct {
	Location        string `json:"location"`
	HAN             string `json:"HAN"`
	NeutralHomeTeam *bool  `json:"neutral_hometeam"`
	Facility        string `json:"facility"`
}

// Links structure
type Links struct {
	Livestats     string      `json:"livestats"`
	LivestatsText string      `json:"livestats_text"`
	Video         string      `json:"video"`
	VideoText     string      `json:"video_text"`
	Audio         string      `json:"audio"`
	AudioText     string      `json:"audio_text"`
	Notes         string      `json:"notes"`
	Tickets       string      `json:"tickets"`
	History       string      `json:"history"`
	BoxScore      *BoxScore   `json:"boxscore"`
	PostGame      *GameInfo   `json:"postgame"`
	PreGame       *GameInfo   `json:"pregame"`
	GameFiles     *[]GameFile `json:"gamefiles"`
}

// BoxScore structure
type BoxScore struct {
	Bid  string `json:"bid"`
	URL  string `json:"url"`
	Text string `json:"text"`
}

// GameInfo structure
type GameInfo struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	StoryImageURL string `json:"story_image_url"`
	Text          string `json:"text"`
	RedirectURL   string `json:"redirect_absolute_url"`
}

// Opponent structure
type Opponent struct {
	GlobalID        int    `json:"opponent_global_id"`
	Name            string `json:"name"`
	Logo            string `json:"logo"`
	LogoImage       string `json:"logo_image"`
	Location        string `json:"location"`
	Mascot          string `json:"mascot"`
	Website         string `json:"opponent_website"`
	ConferenceGame  string `json:"conference_game"`
	Tournament      string `json:"tournament"`
	TournamentColor string `json:"tournament_color"`
}

// Result structure
type Result struct {
	Game           string `json:"game"`
	Status         string `json:"status"`
	TeamScore      string `json:"team_score"`
	OpponentScore  string `json:"opponent_score"`
	PreScoreInfo   string `json:"prescore_info"`
	PostScoreInfo  string `json:"postscore_info"`
	InProgressInfo string `json:"inprogress_info"`
}

// GameFile structure
type GameFile struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

// AllAccessVideo structure
type AllAccessVideo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Format    string `json:"format"`
	Free      bool   `json:"free"`
	URL       string `json:"url"`
	PosterURL string `json:"poster_url"`
}

// SportSeasons structure
type SportSeasons struct {
	Code     string   `json:"code"`
	Label    string   `json:"label"`
	StaffURL string   `json:"staff"`
	Seasons  []Season `json:"schedules"`
}

// Season structure
type Season struct {
	Year         string `json:"year"`
	ScheduleYear string `json:"schedule_year"`
	RosterURL    string `json:"roster"`
	RosterID     *int   `json:"roster_id"`
	ScheduleURL  string `json:"schedule"`
	Current      bool   `json:"current"`
}

// Record structure
type Record struct {
	OverallPercentage    string `json:"overall_percentage"`
	OverallRecord        string `json:"overall_record_unformatted"`
	ConferenceRecord     string `json:"conference_record"`
	ConferencePercentage string `json:"conference_percentage"`
	ConferencePoints     string `json:"conference_points"`
	Streak               string `json:"streak"`
	HomeRecord           string `json:"home_record"`
	AwayRecord           string `json:"away_record"`
	NeutralRecord        string `json:"neutral_record"`
}

// LiveGameItem structure
type LiveGameItem struct {
	GameID       string
	Time         time.Time
	Sport        string
	Home         bool
	OpponentName string
}

// GameItems structure
type GameItems struct {
	Games []model.LiveGame
}
