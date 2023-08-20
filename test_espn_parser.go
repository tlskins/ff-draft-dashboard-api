package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
)

func main() {
	client := p.NewHttpClient()

	players, err := p.GetEspnPlayersForYear(client, 2023)
	if err != nil {
		log.Fatal(err)
	}

	for _, player := range players {
		if player.Id == 4047646 {
			spew.Dump(player)
		}
	}
}
