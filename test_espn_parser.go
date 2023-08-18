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
	spew.Dump(players[0])
}
