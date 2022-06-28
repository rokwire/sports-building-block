package source

// Config structure
type Config struct {
	LivestatsSource    map[string]map[string][]string `json:"livestats_source"`
	FootballConfig     FootballConfig                 `json:"football_config"`
	MBasketballConfig  MBasketballConfig              `json:"mbball_config"`
	WBasketballConfig  WBasketballConfig              `json:"wbball_config"`
	VolleyballConfig   VolleyballConfig               `json:"wvball_config"`
	NotificationConfig NotificationConfig             `json:"notification_config"`
}

// NotificationConfig structure
type NotificationConfig struct {
	Messages map[string]string `json:"messages"`
}

// FootballConfig structure
type FootballConfig struct {
	Phases          map[string]string `json:"phases"`
	LastPlayEnabled bool              `json:"last_play_enabled"`
	XMLDateCheck    bool              `json:"xml_date_check"`
}

// MBasketballConfig structure
type MBasketballConfig struct {
	Phases          map[string]string `json:"phases"`
	LastPlayEnabled bool              `json:"last_play_enabled"`
	XMLDateCheck    bool              `json:"xml_date_check"`
}

// WBasketballConfig structure
type WBasketballConfig struct {
	Phases          map[string]string `json:"phases"`
	LastPlayEnabled bool              `json:"last_play_enabled"`
	XMLDateCheck    bool              `json:"xml_date_check"`
}

// VolleyballConfig structure
type VolleyballConfig struct {
	Phases       map[string]string `json:"phases"`
	XMLDateCheck bool              `json:"xml_date_check"`
}

// NewConfig creates Config instance
func NewConfig() Config {
	var config Config

	//livestats source
	livestatsSource := make(map[string]map[string][]string)

	football := make(map[string][]string)
	football["home"] = []string{"xml_feed", "sidearm"}
	football["away"] = []string{"sidearm"}
	livestatsSource["football"] = football

	mbball := make(map[string][]string)
	mbball["home"] = []string{"xml_feed", "sidearm"}
	mbball["away"] = []string{"sidearm"}
	livestatsSource["mbball"] = mbball

	wbball := make(map[string][]string)
	wbball["home"] = []string{"xml_feed", "sidearm"}
	wbball["away"] = []string{"sidearm"}
	livestatsSource["wbball"] = wbball

	wvball := make(map[string][]string)
	wvball["home"] = []string{"xml_feed", "sidearm"}
	wvball["away"] = []string{"sidearm"}
	livestatsSource["wvball"] = wvball

	mten := make(map[string][]string)
	mten["home"] = []string{"sidearm"}
	mten["away"] = []string{"sidearm"}
	livestatsSource["mten"] = mten

	wten := make(map[string][]string)
	wten["home"] = []string{"sidearm"}
	wten["away"] = []string{"sidearm"}
	livestatsSource["wten"] = wten

	baseball := make(map[string][]string)
	baseball["home"] = []string{"sidearm"}
	baseball["away"] = []string{"sidearm"}
	livestatsSource["baseball"] = baseball

	softball := make(map[string][]string)
	softball["home"] = []string{"sidearm"}
	softball["away"] = []string{"sidearm"}
	livestatsSource["softball"] = softball

	wsoc := make(map[string][]string)
	wsoc["home"] = []string{"sidearm"}
	wsoc["away"] = []string{"sidearm"}
	livestatsSource["wsoc"] = wsoc

	config.LivestatsSource = livestatsSource

	config.FootballConfig = createFootballConfig()
	config.MBasketballConfig = createMBasketballConfig()
	config.WBasketballConfig = createWBasketballConfig()
	config.VolleyballConfig = createVolleyballConfig()
	config.NotificationConfig = createNotificationConfig()

	return config
}

//GetLivestatsSource gives the livestats source
func (config *Config) GetLivestatsSource(sport string, home bool) []string {
	sportConfig := config.LivestatsSource[sport]
	loc := "away"
	if home {
		loc = "home"
	}
	sources := sportConfig[loc]
	return sources
}

//GetFootballDateCheck gives the xml feed date check flag
func (config *Config) GetFootballDateCheck() bool {
	return config.FootballConfig.XMLDateCheck
}

//GetFootballPhaseLabel gives the football phase label
func (config *Config) GetFootballPhaseLabel(phase string) string {
	footbalConfig := config.FootballConfig
	phases := footbalConfig.Phases
	return phases[phase]
}

//GetMBasketballPhaseLabel gives the men's basketball phase label
func (config *Config) GetMBasketballPhaseLabel(phase string) string {
	mbballConfig := config.MBasketballConfig
	phases := mbballConfig.Phases
	return phases[phase]
}

//GetWBasketballPhaseLabel gives the women's basketball phase label
func (config *Config) GetWBasketballPhaseLabel(phase string) string {
	wbballConfig := config.WBasketballConfig
	phases := wbballConfig.Phases
	return phases[phase]
}

//GetBasketballDateCheck gives the basketball date check flag
func (config *Config) GetBasketballDateCheck(sport string) bool {
	switch sport {
	case "mbball":
		return config.MBasketballConfig.XMLDateCheck
	case "wbball":
		return config.WBasketballConfig.XMLDateCheck
	default:
		return false
	}
}

//GetBasketballLastPlay gives the basketball last play flag
func (config *Config) GetBasketballLastPlay(sport string) bool {
	switch sport {
	case "mbball":
		return config.MBasketballConfig.LastPlayEnabled
	case "wbball":
		return config.WBasketballConfig.LastPlayEnabled
	default:
		return false
	}
}

//GetVolleyballDateCheck gives the volleyball date check flag
func (config *Config) GetVolleyballDateCheck() bool {
	return config.VolleyballConfig.XMLDateCheck
}

//GetVolleyballPhaseLabel gives the volleyball phase label
func (config *Config) GetVolleyballPhaseLabel(phase string) string {
	wvballConfig := config.VolleyballConfig
	phases := wvballConfig.Phases
	return phases[phase]
}

func createFootballConfig() FootballConfig {
	var footballConfig FootballConfig

	phases := make(map[string]string)
	phases["pre"] = "Pregame"
	phases["1"] = "1st Quarter"
	phases["2"] = "2nd Quarter"
	phases["ht"] = "Half Time"
	phases["3"] = "3rd Quarter"
	phases["4"] = "4th Quarter"
	phases["ot"] = "Over Time"
	phases["final"] = "Final Score"
	footballConfig.Phases = phases

	footballConfig.LastPlayEnabled = true
	footballConfig.XMLDateCheck = true

	return footballConfig
}

func createMBasketballConfig() MBasketballConfig {
	var mBasketballConfig MBasketballConfig

	phases := make(map[string]string)
	phases["pre"] = "Pregame"
	phases["1"] = "1st Half"
	phases["2"] = "2nd Half"
	phases["ot"] = "Over Time"
	phases["final"] = "Final Score"
	mBasketballConfig.Phases = phases

	mBasketballConfig.LastPlayEnabled = false
	mBasketballConfig.XMLDateCheck = true

	return mBasketballConfig
}

func createWBasketballConfig() WBasketballConfig {
	var wBasketballConfig WBasketballConfig

	phases := make(map[string]string)
	phases["pre"] = "Pregame"
	phases["1"] = "1st Quarter"
	phases["2"] = "2nd Quarter"
	phases["3"] = "3rd Quarter"
	phases["4"] = "4th Quarter"
	phases["ot"] = "Over Time"
	phases["final"] = "Final Score"
	wBasketballConfig.Phases = phases

	wBasketballConfig.LastPlayEnabled = false
	wBasketballConfig.XMLDateCheck = true

	return wBasketballConfig
}

func createVolleyballConfig() VolleyballConfig {
	var volleyballConfig VolleyballConfig

	phases := make(map[string]string)
	phases["pre"] = "Pregame"
	phases["game_name"] = "Set"
	phases["final"] = "Final Score"
	volleyballConfig.Phases = phases

	volleyballConfig.XMLDateCheck = true

	return volleyballConfig
}

func createNotificationConfig() NotificationConfig {
	var notificationConfig NotificationConfig

	messages := make(map[string]string)
	messages["game_title_format"] = "%[1]s vs %[2]s" // [1] - home team name, [2] - away team name
	messages["game_started_msg"] = "The Game has started"
	messages["game_ended_msg"] = "The Game had ended."
	messages["game_ended_score_format"] = "Score %[1]s %[2]d : %[3]s %[4]d" // [1] - home team name, [2] - home score, [3] - away team name, [4] - away score
	messages["news_updates_sport_title_format"] = "Athletics news - %[1]s"  // [1] - news sport / category
	messages["news_updates_default_title"] = "Athletics news"
	messages["news_updates_body_content_format"] = "%[1]s" //[1] - news title
	notificationConfig.Messages = messages

	return notificationConfig
}
