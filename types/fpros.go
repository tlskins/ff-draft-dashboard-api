package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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
	Id                   int     `json:"player_id"`
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

	// calculated
	PosStdRank string `json:"posStdRank"`
	PosPprRank string `json:"posPprRank"`
	StdRank    int    `json:"stdRank"`
	PprRank    int    `json:"pprRank"`
}

func (p FproPlayer) Team() string {
	if p.PlayerTeamId == "PHI" {
		return "PHL"
	} else {
		return p.PlayerTeamId
	}
}

func (p FproPlayer) Position() Position {
	return Position(p.PlayerPositionId)
}

func (p FproPlayer) CleanName() string {
	alphaRgx := regexp.MustCompile("[^a-zA-Z]+")
	nameKey := strings.ToUpper(p.PlayerName)
	nameKey, _ = strings.CutSuffix(nameKey, " JR.")
	nameKey, _ = strings.CutSuffix(nameKey, " III")
	nameKey, _ = strings.CutSuffix(nameKey, " II")
	nameKey = alphaRgx.ReplaceAllString(nameKey, "")

	return nameKey
}

func (p FproPlayer) PrimaryKey() string {
	return fmt.Sprintf("%s-%s-%s", p.CleanName(), strings.ToUpper(p.Team()), p.Position())
}

func (p FproPlayer) PosRankInt(isPpr bool) int {
	reg := regexp.MustCompile("[0-9]+")
	posRank := p.PosPprRank
	if !isPpr {
		posRank = p.PosStdRank
	}
	match := reg.FindString(posRank)
	if match != "" {
		num, _ := strconv.Atoi(match)
		return num
	}

	return 0
}

func (p FproPlayer) ToPlayer() (out *Player) {
	names := strings.Split(p.PlayerName, " ")
	out = &Player{
		Id:            p.SportsDataId,
		FirstName:     names[0],
		LastName:      strings.Join(names[1:], " "),
		Name:          p.PlayerName,
		MatchName:     MatchName(p.PlayerName),
		Position:      p.Position(),
		Team:          p.Team(),
		CustomStdRank: p.RankEcr,
		CustomPprRank: p.RankEcr,
		Tier:          strconv.Itoa(p.Tier),
	}
	return
}

type FprosPlayerMatch struct {
	Fpros  *FproPlayer `json:"fprosPlayer"`
	Player *Player     `json:"player"`
}

func (m *FprosPlayerMatch) AddPlayerRank() {
	if m.Player != nil && m.Fpros != nil {
		m.Player.CustomStdOvrRank = m.Fpros.StdRank
		m.Player.CustomPprOvrRank = m.Fpros.PprRank
		m.Player.CustomStdRank = m.Fpros.PosRankInt(false)
		m.Player.CustomPprRank = m.Fpros.PosRankInt(true)
		m.Player.Tier = strconv.Itoa(m.Fpros.Tier)
	}

	return
}
