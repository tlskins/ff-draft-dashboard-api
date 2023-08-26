package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
)

func main() {
	client := p.NewHttpClient()
	fprosPlayers, err := p.GetFprosPlayers(client)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(fprosPlayers[0].PosRankInt(true))
}
