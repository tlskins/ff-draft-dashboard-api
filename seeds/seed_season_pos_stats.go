package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func main() {
	// vars
	currYear := 2023
	numTeams := 12

	client := p.NewHttpClient()
	espnPlayers, err := p.GetEspnPlayersForYear(client, currYear, 350)
	if err != nil {
		log.Fatal(err)
	}

	players := make([]*t.Player, len(espnPlayers))
	for i, espnPlayer := range espnPlayers {
		players[i] = espnPlayer.ToPlayer()
	}

	harrisPlayers := p.ParseHarrisRanksV2(currYear)
	fmt.Printf("%v harris players\n", len(harrisPlayers))

	_, err = p.AddHarrisRanks(harrisPlayers, players)
	if err != nil {
		panic(err)
	}

	yrStatsByPos := t.CalcStatsForYear(players, numTeams, currYear-1)

	spew.Dump(yrStatsByPos)
}
