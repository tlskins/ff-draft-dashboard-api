package main

import (
	"fmt"
	"log"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
)

func main() {
	client := p.NewHttpClient()
	espnPlayers, err := p.GetEspnPlayersForYear(client, 2023)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v espn players\n", len(espnPlayers))
	harrisPlayers := p.ParseHarrisRanksV2(2023)
	fmt.Printf("%v harris players\n", len(harrisPlayers))

	players, unmatched, err := p.MatchHarrisAndEspnPlayers(harrisPlayers, espnPlayers)
	if err != nil {
		panic(err)
	}
	for _, player := range unmatched {
		log.Printf("Unfound: %s %s %s: %v %v\n", player.Name, player.Position, player.Team, player.EspnOvrStdRank, player.CustomStdRank)
	}

	fmt.Printf("total players %v\n", len(players))
	fmt.Printf("unmatched %v\n", len(unmatched))
}
