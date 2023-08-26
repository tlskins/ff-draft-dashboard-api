package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly/v2"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func main() {

	ParsePlayersForPath := func(pos, urlPath string, isRankByType bool) (players []*t.HarrisPlayer) {
		players = []*t.HarrisPlayer{}
		playersMap := map[string]*t.HarrisPlayer{}
		c := colly.NewCollector()

		c.OnHTML("body", func(bodyEl *colly.HTMLElement) {
			bodyEl.ForEach("table", func(idx int, tableEl *colly.HTMLElement) {
				var isPpr *int // -1 no, 1 yes, 0 both
				if !isRankByType {
					val := 0
					isPpr = &val
				}

				tableEl.ForEach("tr", func(idx int, rowEl *colly.HTMLElement) {
					newPlayer := t.HarrisPlayer{Position: pos}
					playerFound := false
					var rank int

					rowEl.ForEach("td", func(idx int, tdEl *colly.HTMLElement) {
						cellText := tdEl.Text

						// determine ranking type
						if isPpr == nil && isRankByType {
							if strings.Contains(strings.ToLower(cellText), "standard") {
								val := -1
								isPpr = &val
							} else if strings.Contains(strings.ToLower(cellText), "ppr") {
								val := 1
								isPpr = &val
							}
						} else if isPpr != nil { // if ranking type is determined start tracking players
							if idx == 0 {
								// ranking
								var rankErr error
								rank, rankErr = strconv.Atoi(cellText)
								if rankErr == nil {
									playerFound = true
								}
							} else if idx == 1 {
								// name
								newPlayer.Name = cellText
								newPlayer.FirstName, newPlayer.LastName = p.ParseHarrisName(newPlayer.Name)
							} else if idx == 2 {
								// team
								newPlayer.Team = cellText
							}
						}
					})

					// player parsing complete
					// fmt.Sprintf("found %v, ppr %v, rank %v\n", playerFound, isPpr, rank)
					if playerFound && isPpr != nil && rank > 0 {
						newPlayer.Id = p.HarrisPlayerPrimaryKey(&newPlayer)

						// get player to update
						currPlayer := &newPlayer
						existPlayer := playersMap[newPlayer.Id]
						if existPlayer != nil {
							currPlayer = existPlayer
						}

						// set rank
						if *isPpr == 1 || *isPpr == 0 {
							currPlayer.PPRRank = rank
						}
						if *isPpr == -1 || *isPpr == 0 {
							currPlayer.StdRank = rank
						}

						// is a new player add to trackers
						if existPlayer == nil {
							playersMap[newPlayer.Id] = &newPlayer
							players = append(players, &newPlayer)
						}
					}
				})
			})
		})
		c.Visit(fmt.Sprintf("https://www.harrisfootball.com/%s", urlPath))

		return
	}

	posToPathMap := map[string]string{
		"QB": "ranks-draft",
		"RB": "rb-ranks-draft",
		"WR": "wr-ranks-draft",
		"TE": "te-ranks-draft",
	}
	posToRankTypeMap := map[string]bool{
		"QB": false,
		"RB": true,
		"WR": true,
		"TE": false,
	}
	allPlayers := []*t.HarrisPlayer{}
	for pos, urlPath := range posToPathMap {
		posPlayers := ParsePlayersForPath(pos, urlPath, posToRankTypeMap[pos])
		allPlayers = append(allPlayers, posPlayers...)
		spew.Dump(posPlayers[0])
		spew.Dump(posPlayers[len(posPlayers)-1])
	}
	fmt.Printf("All Players %v\n", len(allPlayers))
}
