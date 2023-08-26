package parsers

import (
	"math"
	"sort"

	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func GetPlayersStats(year, stIdx, endIdx int, players []*t.Player) (out []*t.Stats) {
	out = []*t.Stats{}
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

// run on front end instead
func CalcStatsForPosPlayers(posPlayers []*t.Player, ssnStats *t.SeasonPositionalStats, numTeams int) {
	for _, player := range posPlayers {
		player.CalcStats(ssnStats, numTeams)
	}
}

// creates cumulative stats by pos / year
func CalcStatsForYear(players []*t.Player, numTeams, year int) (
	yrStatsByPos map[t.Position]*t.SeasonPositionalStats,
	qbs, rbs, wrs, tes []*t.Player,
) {
	qbs, rbs, wrs, tes = GroupPlayersForYear(players, year)

	yrStatsByPos = map[t.Position]*t.SeasonPositionalStats{
		t.QB: CalcAvgStatsForPos(qbs, numTeams, year, t.QB),
		t.RB: CalcAvgStatsForPos(rbs, numTeams, year, t.RB),
		t.WR: CalcAvgStatsForPos(wrs, numTeams, year, t.WR),
		t.TE: CalcAvgStatsForPos(tes, numTeams, year, t.TE),
	}

	return
}

// mutates players and adds last yr tier to players and
func AddStatsForPosPlayers(qbs, rbs, wrs, tes []*t.Player, yrStatsByPos map[t.Position]*t.SeasonPositionalStats, numTeams int) {
	// calc player stats
	CalcStatsForPosPlayers(qbs, yrStatsByPos[t.QB], numTeams)
	CalcStatsForPosPlayers(rbs, yrStatsByPos[t.RB], numTeams)
	CalcStatsForPosPlayers(wrs, yrStatsByPos[t.WR], numTeams)
	CalcStatsForPosPlayers(tes, yrStatsByPos[t.TE], numTeams)
}

func GroupPlayersForYear(players []*t.Player, year int) (qbs, rbs, wrs, tes []*t.Player) {
	qbs = []*t.Player{}
	rbs = []*t.Player{}
	wrs = []*t.Player{}
	tes = []*t.Player{}
	for _, player := range players {
		if player.Position == t.QB {
			qbs = append(qbs, player)
		} else if player.Position == t.RB {
			rbs = append(rbs, player)
		} else if player.Position == t.WR {
			wrs = append(wrs, player)
		} else if player.Position == t.TE {
			tes = append(tes, player)
		}
	}

	sort.Slice(qbs, func(i, j int) bool {
		return qbs[i].TotalPtsForYear(year) > qbs[j].TotalPtsForYear(year)
	})
	sort.Slice(rbs, func(i, j int) bool {
		return rbs[i].TotalPtsForYear(year) > rbs[j].TotalPtsForYear(year)
	})
	sort.Slice(wrs, func(i, j int) bool {
		return wrs[i].TotalPtsForYear(year) > wrs[j].TotalPtsForYear(year)
	})
	sort.Slice(tes, func(i, j int) bool {
		return tes[i].TotalPtsForYear(year) > tes[j].TotalPtsForYear(year)
	})

	return
}

// pos players should be sorted by fps for that year
func CalcAvgStatsForPos(posPlayers []*t.Player, numTeams, year int, pos t.Position) (out *t.SeasonPositionalStats) {
	out = &t.SeasonPositionalStats{
		Id:       t.SeasonPositionalRankId(pos, year),
		Year:     year,
		Position: pos,
	}

	if pos == t.WR || pos == t.RB {
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

func AvgStats(allStats []*t.Stats) (out *t.Stats) {
	out = &t.Stats{}
	if len(allStats) == 0 {
		return
	}
	count := float64(len(allStats))
	minPPG := math.MaxFloat64
	maxPPG := 0.0
	ttlPPG := 0.0
	minTtlPts := math.MaxFloat64
	maxTtlPts := 0.0
	ttlTtlPts := 0.0
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
		if minTtlPts > stats.TotalPoints {
			minTtlPts = stats.TotalPoints
		}
		if maxTtlPts < stats.TotalPoints {
			maxTtlPts = stats.TotalPoints
		}
		ttlTtlPts += stats.TotalPoints
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
	out.TotalPoints = ttlTtlPts / count
	out.MinTtlPts = minTtlPts
	out.MaxTtlPts = maxTtlPts
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
