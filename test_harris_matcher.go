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

	matchedPlayers := p.MatchHarrisAndEspnPlayers(harrisPlayers, espnPlayers)
	unmatched := 0
	for _, match := range matchedPlayers {
		player, err := match.ToPlayer()
		if err != nil {
			log.Fatal(err)
		}
		if match.Harris == nil || match.Espn == nil {
			unmatched += 1
			log.Printf("Unfound: %s %s %s: %v %v\n", player.Name, player.Position, player.Team, player.EspnOvrStdRank, player.CustomStdRank)
		}
	}

	fmt.Printf("unmatched %v\n", unmatched)
}
