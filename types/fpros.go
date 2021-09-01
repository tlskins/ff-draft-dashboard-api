package types

type FproEcrData struct {
	Sport        string        `json:"sport"`
	Year         string        `json:"year"`
	Week         string        `json:"week"`
	PositionId   string        `json:"position_id"`
	Scoring      string        `json:"scoring"`
	Count        int           `json:"count"`
	TotalExperts int           `json:"total_experts"`
	LastUpdated  string        `json:"last_updated"`
	Players      []*FproPlayer `json:"players"`
}

type FproPlayer struct {
	RankEcr              int     `json:"rank_ecr"`
	SportsDataId         string  `json:"sportsdata_id"`
	PlayerPageUrl        string  `json:"player_page_url"`
	PlayerSquareImageUrl string  `json:"player_square_image_url"`
	PlayerOwnedYahoo     float64 `json:"player_owned_yahoo"`
	RankMin              string  `json:"rank_min"`
	RankMax              string  `json:"rank_max"`
	PlayerName           string  `json:"player_name"`
	PlayerOwnedAvg       float64 `json:"player_owned_avg"`
	PlayerFilename       string  `json:"player_filename"`
	PlayerShortName      string  `json:"player_short_name"`
	PlayerEligibility    string  `json:"player_eligibility"`
	RankStd              string  `json:"rank_std"`
	Tier                 int     `json:"tier"`
	PlayerByeWeek        string  `json:"player_bye_week"`
	PosRank              string  `json:"pos_rank"`
	PlayerTeamId         string  `json:"player_team_id"`
	PlayerPositionId     string  `json:"player_position_id"`
}
