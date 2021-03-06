package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func main() {
	client := p.NewHttpClient()
	out := t.EspnPlayersResp{}
	if err := p.HttpRequest(client, "GET", p.EspnApiUrl, p.EspnQueryHeader(350, 0), nil, &out); err != nil {
		log.Fatal(err)
	}
	players := make([]*t.Player, len(out.Players))
	for i, p := range out.Players {
		if p.Profile != nil && p.Profile.FullName == "Sony Michel" {
			spew.Dump(p)
		}
		players[i] = p.ToPlayer()
	}

	currId := 1
	// players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/ranks-draft", t.QB, currId, players)
	// players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/wr-ranks-draft", t.WR, currId, players)
	players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/rb-ranks-draft", t.RB, currId, players)
	// players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/te-ranks-draft", t.TE, currId, players)

	fmt.Println(len(players))
	for _, player := range players {
		if t.MatchName("Sony Michel") == player.MatchName {
			spew.Dump(player)
		}
	}
}
