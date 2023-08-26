package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
)

func main() {
	currYear := 2023

	client := p.NewHttpClient()
	players, err := p.GetEspnPlayersForYear(client, currYear, 350)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v players\n", len(players))

	fprosPlayers, err := p.GetFprosPlayers(client)
	if err != nil {
		log.Fatal(err)
	}

	unmatched := p.AddFprosRanks(fprosPlayers, players)
	for _, player := range unmatched {
		log.Printf("Unfound: %s %s %s: %v %v\n", player.Name, player.Position, player.Team, player.EspnOvrStdRank, player.CustomStdRank)
	}

	fmt.Printf("total players %v\n", len(players))
	fmt.Printf("unmatched %v\n", len(unmatched))

	spew.Dump(players[10])
}
