package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Position string

const (
	QB         Position = "QB"
	RB         Position = "RB"
	WR         Position = "WR"
	TE         Position = "TE"
	DST        Position = "DST"
	NoPosition Position = ""
)

type Player struct {
	Id        string   `json:"id"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Name      string   `json:"name"`
	MatchName string   `json:"matchName"`
	Position  Position `json:"position"`
	Team      string   `json:"team"`
	Tier      string   `json:"tier"`

	CustomStdRank     int     `json:"customStdRank,omitempty"`
	CustomPprRank     int     `json:"customPprRank,omitempty"`
	EspnOvrPprRank    int     `json:"espnOvrPprRank,omitempty"`
	EspnOvrStdRank    int     `json:"espnOvrStdRank,omitempty"`
	EspnAdp           float64 `json:"espnAdp,omitempty"`
	EspnPlayerOutlook string  `json:"espnPlayerOutlook,omitempty"`

	// stats
	SeasonStats   []*Stats `json:"seasonStats"`
	LastYrTier    float64  `json:"lastYrTier"`    // calculated front end based on num teams
	LastYrOvrRank int      `json:"lastYrOvrRank"` // calculated front end based on num teams
	StdRankTier   int      `json:"stdRankTier"`
	PprRankTier   int      `json:"pprRankTier"`
	Pros          string   `json:"pros"`
	Cons          string   `json:"cons"`
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// do this on the front end instead bc its based on num of teams
func (p *Player) CalcStats(lastSsnPosStats *SeasonPositionalStats, numTeams int) {
	if p.Position == RB || p.Position == WR {
		p.StdRankTier = (p.CustomStdRank / numTeams) + 1
		p.PprRankTier = (p.CustomPprRank / numTeams) + 1
	} else {
		p.StdRankTier = MinInt(p.CustomStdRank, numTeams)
		p.PprRankTier = MinInt(p.CustomPprRank, numTeams)
	}
	year := lastSsnPosStats.Year
	var yearStats *Stats
	for _, stats := range p.SeasonStats {
		if stats.Year == year {
			yearStats = stats
		}
	}
	if yearStats == nil {
		return
	}
	if yearStats.TotalPoints >= lastSsnPosStats.Tier1Stats.MinTtlPts {
		p.LastYrTier = 1
	} else if yearStats.TotalPoints >= lastSsnPosStats.Tier2Stats.MinTtlPts {
		p.LastYrTier = 2
	} else if yearStats.TotalPoints >= lastSsnPosStats.Tier3Stats.MinTtlPts {
		p.LastYrTier = 3
	} else if yearStats.TotalPoints >= lastSsnPosStats.Tier4Stats.MinTtlPts {
		p.LastYrTier = 4
	} else if yearStats.TotalPoints >= lastSsnPosStats.Tier5Stats.MinTtlPts {
		p.LastYrTier = 5
	} else {
		p.LastYrTier = 6
	}
}

func (p Player) StatsForYear(year int) *Stats {
	if p.SeasonStats == nil {
		return nil
	}
	for _, stats := range p.SeasonStats {
		if stats.Year == year {
			return stats
		}
	}
	return nil
}

func (p Player) PPGForYear(year int) float64 {
	stats := p.StatsForYear(year)
	if stats == nil {
		return 0
	}
	return stats.PPG
}

func (p Player) TotalPtsForYear(year int) float64 {
	stats := p.StatsForYear(year)
	if stats == nil {
		return 0
	}
	return stats.TotalPoints
}

func (p *Player) AddPlayerReport(report *PlayerReport) {
	p.Pros = report.Pros
	p.Cons = report.Cons
}

// type CustomPlayerRanks interface {
// 	StdRank() int
// 	PprRank() int
// }

// func (p *Player) AddCustomRanks(customPlayerRank CustomPlayerRanks) {
// 	p.CustomStdRank = customPlayerRank.StdRank()
// 	p.CustomPprRank = customPlayerRank.PprRank()
// }

func CleanName(name string) (out string) {
	cleanNmRgx := regexp.MustCompile(`\(.*\)`)
	out = strings.TrimSpace(cleanNmRgx.ReplaceAllString(name, ""))
	return
}

func MatchName(name string) (out string) {
	out = strings.ToLower(name)
	cleanRgx := regexp.MustCompile(`[^a-z ]`)
	out = cleanRgx.ReplaceAllString(out, "")
	out = strings.TrimSpace(out)
	return
}

func FindPlayer(players []*Player, matchName string) (out *Player) {
	matchParts := strings.Split(matchName, " ")
	matchFirst := matchParts[0]
	matchLast := strings.Join(matchParts[1:], "")
	for _, player := range players {
		nameParts := strings.Split(player.MatchName, " ")
		nameFirst := nameParts[0]
		nameLast := strings.Join(nameParts[1:], "")
		if (len(matchFirst) <= len(nameFirst) && !regexp.MustCompile(matchFirst).MatchString(nameFirst)) ||
			(len(nameFirst) < len(matchFirst) && !regexp.MustCompile(nameFirst).MatchString(matchFirst)) {
			continue
		}
		if (len(matchLast) <= len(nameLast) && !regexp.MustCompile(matchLast).MatchString(nameLast)) ||
			(len(nameLast) < len(matchLast) && !regexp.MustCompile(nameLast).MatchString(matchLast)) {
			continue
		}
		return player
	}
	return
}

type Stats struct {
	Year            int     `bson:"year" json:"year"`
	TotalPoints     float64 `bson:"totalPoints" json:"totalPoints"`
	MinTtlPts       float64 `bson:"minTotalPts" json:"minTotalPts"`
	MaxTtlPts       float64 `bson:"maxTotalPts" json:"maxTotalPts"`
	PPG             float64 `bson:"ppg" json:"ppg"`
	MinPPG          float64 `bson:"minPPG,omitempty" json:"minPPG,omitempty"`
	MaxPPG          float64 `bson:"maxPPG,omitempty" json:"maxPPG,omitempty"`
	GamesPlayed     float64 `bson:"gamesPlayed" json:"gamesPlayed"`
	RushAttempts    float64 `bson:"rushAttempts" json:"rushAttempts"`
	RushYards       float64 `bson:"rushYards" json:"rushYards"`
	RushTds         float64 `bson:"rushTds" json:"rushTds"`
	Recs            float64 `bson:"recs" json:"recs"`
	RecYards        float64 `bson:"recYards" json:"recYards"`
	RecTds          float64 `bson:"recTds" json:"recTds"`
	PassAttempts    float64 `bson:"passAttempts" json:"passAttempts"`
	PassCompletions float64 `bson:"passCompletions" json:"passCompletions"`
	PassYards       float64 `bson:"passYards" json:"passYards"`
	PassTds         float64 `bson:"passTds" json:"passTds"`
	PassInts        float64 `bson:"passInts" json:"passInts"`
}

func SeasonPositionalRankId(pos Position, year int) string {
	return fmt.Sprintf("%v-%s", year, pos)
}

type SeasonPositionalStats struct {
	Id         string   `bson:"_id" json:"id"`
	Year       int      `bson:"year" json:"year"`
	Position   Position `bson:"pos" json:"pos"`
	Tier1Stats *Stats   `bson:"tier1Stats" json:"tier1Stats"`
	Tier2Stats *Stats   `bson:"tier2Stats" json:"tier2Stats"`
	Tier3Stats *Stats   `bson:"tier3Stats" json:"tier3Stats"`
	Tier4Stats *Stats   `bson:"tier4Stats" json:"tier4Stats"`
	Tier5Stats *Stats   `bson:"tier5Stats" json:"tier5Stats"`
	Tier6Stats *Stats   `bson:"tier6Stats" json:"tier6Stats"`
}

type PlayerReport struct {
	Id        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	Name      string    `bson:"nm" json:"name"`
	Position  Position  `bson:"pos" json:"position"`
	Pros      string    `bson:"pros" json:"pros"`
	Cons      string    `bson:"cons" json:"cons"`
}
