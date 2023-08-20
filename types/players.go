package types

import (
	"fmt"
	"math"
	"regexp"
	"sort"
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
	SeasonStats []*Stats `json:"seasonStats"`
	LastYrTier  float64  `json:"lastYrTier"`
	StdRankTier int      `json:"stdRankTier"`
	PprRankTier int      `json:"pprRankTier"`
	Pros        string   `json:"pros"`
	Cons        string   `json:"cons"`
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

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
	if yearStats.PPG >= lastSsnPosStats.Tier1Stats.MinPPG {
		p.LastYrTier = 1
	} else if yearStats.PPG >= lastSsnPosStats.Tier2Stats.MinPPG {
		p.LastYrTier = 2
	} else if yearStats.PPG >= lastSsnPosStats.Tier3Stats.MinPPG {
		p.LastYrTier = 3
	} else if yearStats.PPG >= lastSsnPosStats.Tier4Stats.MinPPG {
		p.LastYrTier = 4
	} else if yearStats.PPG >= lastSsnPosStats.Tier5Stats.MinPPG {
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

func GetPlayersStats(year, stIdx, endIdx int, players []*Player) (out []*Stats) {
	out = []*Stats{}
	if stIdx > len(players) {
		return
	}
	if endIdx > len(players) {
		endIdx = len(players)
	}
	for _, player := range players[stIdx:endIdx] {
		if player != nil && player.SeasonStats != nil {
			for _, stats := range player.SeasonStats {
				if stats.Year == year {
					out = append(out, stats)
				}
			}
		}
	}

	return
}

func CalcStatsForPosPlayers(posPlayers []*Player, ssnStats *SeasonPositionalStats, numTeams int) {
	for _, player := range posPlayers {
		player.CalcStats(ssnStats, numTeams)
	}
}

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

func AvgStats(allStats []*Stats) (out *Stats) {
	out = &Stats{}
	if len(allStats) == 0 {
		return
	}
	count := float64(len(allStats))
	minPPG := math.MaxFloat64
	maxPPG := 0.0
	ttlPPG := 0.0
	ttlGamesPlayed := 0.0
	ttlRushAtts := 0.0
	ttlRushYds := 0.0
	ttlRushTds := 0.0
	ttlRecs := 0.0
	ttlRecYds := 0.0
	ttlRecTds := 0.0
	ttlPassAtts := 0.0
	ttlPassComps := 0.0
	ttlPassYds := 0.0
	ttlPassTds := 0.0
	ttlPassInts := 0.0
	for _, stats := range allStats {
		if out.Year == 0 {
			out.Year = stats.Year
		}
		if minPPG > stats.PPG {
			minPPG = stats.PPG
		}
		if maxPPG < stats.PPG {
			maxPPG = stats.PPG
		}
		ttlPPG += stats.PPG
		ttlGamesPlayed += stats.GamesPlayed
		ttlRushAtts += stats.RushAttempts
		ttlRushYds += stats.RushYards
		ttlRushTds += stats.RushTds
		ttlRecs += stats.Recs
		ttlRecYds += stats.RecYards
		ttlRecTds += stats.RecTds
		ttlPassAtts += stats.PassAttempts
		ttlPassComps += stats.PassCompletions
		ttlPassYds += stats.PassYards
		ttlPassTds += stats.PassTds
		ttlPassInts += stats.PassInts
	}
	out.PPG = ttlPPG / count
	out.MinPPG = minPPG
	out.MaxPPG = maxPPG
	out.GamesPlayed = ttlGamesPlayed / count
	out.RushAttempts = ttlRushAtts / count
	out.RushYards = ttlRushYds / count
	out.RushTds = ttlRushTds / count
	out.Recs = ttlRecs / count
	out.RecYards = ttlRecYds / count
	out.RecTds = ttlRecTds / count
	out.PassAttempts = ttlPassAtts / count
	out.PassCompletions = ttlPassComps / count
	out.PassYards = ttlPassYds / count
	out.PassTds = ttlPassTds / count
	out.PassInts = ttlPassInts / count

	return
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

func GroupPlayersForYear(players []*Player, year int) (qbs, rbs, wrs, tes []*Player) {
	qbs = []*Player{}
	rbs = []*Player{}
	wrs = []*Player{}
	tes = []*Player{}
	for _, player := range players {
		if player.Position == QB {
			qbs = append(qbs, player)
		} else if player.Position == RB {
			rbs = append(rbs, player)
		} else if player.Position == WR {
			wrs = append(wrs, player)
		} else if player.Position == TE {
			tes = append(tes, player)
		}
	}

	sort.Slice(qbs, func(i, j int) bool {
		return qbs[i].PPGForYear(year) > qbs[j].PPGForYear(year)
	})
	sort.Slice(rbs, func(i, j int) bool {
		return rbs[i].PPGForYear(year) > rbs[j].PPGForYear(year)
	})
	sort.Slice(wrs, func(i, j int) bool {
		return wrs[i].PPGForYear(year) > wrs[j].PPGForYear(year)
	})
	sort.Slice(tes, func(i, j int) bool {
		return tes[i].PPGForYear(year) > tes[j].PPGForYear(year)
	})

	return
}

// pos players should be sorted by fps for that year
func CalcAvgStatsForPos(posPlayers []*Player, numTeams, year int, pos Position) (out *SeasonPositionalStats) {
	out = &SeasonPositionalStats{
		Id:       SeasonPositionalRankId(pos, year),
		Year:     year,
		Position: pos,
	}

	if pos == WR || pos == RB {
		out.Tier1Stats = AvgStats(GetPlayersStats(year, 0, numTeams+1, posPlayers))
		out.Tier2Stats = AvgStats(GetPlayersStats(year, numTeams+1, numTeams*2+1, posPlayers))
		out.Tier3Stats = AvgStats(GetPlayersStats(year, numTeams*2+1, numTeams*3+1, posPlayers))
		out.Tier4Stats = AvgStats(GetPlayersStats(year, numTeams*3+1, numTeams*4+1, posPlayers))
		out.Tier5Stats = AvgStats(GetPlayersStats(year, numTeams*4+1, numTeams*5+1, posPlayers))
		out.Tier6Stats = AvgStats(GetPlayersStats(year, numTeams*5+1, numTeams*6+1, posPlayers))
	} else {
		out.Tier1Stats = posPlayers[0].StatsForYear(year)
		out.Tier2Stats = posPlayers[1].StatsForYear(year)
		out.Tier3Stats = posPlayers[2].StatsForYear(year)
		out.Tier4Stats = posPlayers[3].StatsForYear(year)
		out.Tier5Stats = posPlayers[4].StatsForYear(year)
		out.Tier6Stats = posPlayers[5].StatsForYear(year)
	}

	return
}

type PlayerReport struct {
	Id        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	Name      string    `bson:"nm" json:"name"`
	Position  Position  `bson:"pos" json:"position"`
	Pros      string    `bson:"pros" json:"pros"`
	Cons      string    `bson:"cons" json:"cons"`
}
