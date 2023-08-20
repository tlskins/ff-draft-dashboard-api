package main

import (
	"fmt"
	"log"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func main() {
	// vars
	currYear := 2023
	numTeams := 12

	client := p.NewHttpClient()
	espnPlayers, err := p.GetEspnPlayersForYear(client, currYear)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v espn players\n", len(espnPlayers))
	harrisPlayers := p.ParseHarrisRanksV2(2023)
	fmt.Printf("%v harris players\n", len(harrisPlayers))
	players, _, err := p.MatchHarrisAndEspnPlayers(harrisPlayers, espnPlayers)
	if err != nil {
		panic(err)
	}

	qbs, rbs, wrs, tes := t.GroupPlayersForYear(players, currYear-1)

	fmt.Printf("Found %v qbs\n", len(qbs))
	fmt.Printf("Found %v rbs\n", len(rbs))
	fmt.Printf("Found %v wrs\n", len(wrs))
	fmt.Printf("Found %v tes\n", len(tes))

	// calculate pos stats
	lastYrStatsByPos := map[t.Position]*t.SeasonPositionalStats{
		t.QB: t.CalcAvgStatsForPos(qbs, numTeams, currYear-1, t.QB),
		t.RB: t.CalcAvgStatsForPos(rbs, numTeams, currYear-1, t.RB),
		t.WR: t.CalcAvgStatsForPos(wrs, numTeams, currYear-1, t.WR),
		t.TE: t.CalcAvgStatsForPos(tes, numTeams, currYear-1, t.TE),
	}

	// calc player stats
	t.CalcStatsForPosPlayers(qbs, lastYrStatsByPos[t.QB], numTeams)
	t.CalcStatsForPosPlayers(rbs, lastYrStatsByPos[t.RB], numTeams)
	t.CalcStatsForPosPlayers(wrs, lastYrStatsByPos[t.WR], numTeams)
	t.CalcStatsForPosPlayers(tes, lastYrStatsByPos[t.TE], numTeams)
}
