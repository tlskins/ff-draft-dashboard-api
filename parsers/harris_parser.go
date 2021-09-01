package parsers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func ParseHarrisRanks(url string, pos t.Position, currId int, players []*t.Player) ([]*t.Player, int) {
	isRankRgx := regexp.MustCompile(`^[0-9]+$`)
	isTmRgx := regexp.MustCompile(`^[A-Z]{2,}$`)
	isStdScrRgx := regexp.MustCompile(`(?i)standard scoring`)
	isPprScrRgx := regexp.MustCompile(`(?i)ppr scoring`)

	rankType := "BOTH" // PPR / STD / BOTH

	c := colly.NewCollector()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		texts := e.ChildTexts("table > tbody > tr > td")
		rank := 1
		team := ""
		name := ""
		matchName := ""
		fName := ""
		lName := ""

		for _, text := range texts {
			if len(text) == 0 {
				continue
			}
			if isRankRgx.MatchString(text) {
				rank, _ = strconv.Atoi(text)
			}
			if isStdScrRgx.MatchString(text) {
				rankType = "STD"
				continue
			}
			if isPprScrRgx.MatchString(text) {
				rankType = "PPR"
				continue
			}

			if !isTmRgx.MatchString(text) {
				name = t.CleanName(text)
				matchName = t.MatchName(text)
				nameParts := strings.Split(name, " ")
				fName = nameParts[0]
				lName = strings.Join(nameParts[1:], " ")
				continue
			} else {
				team = strings.TrimSpace(text)
				player := t.FindPlayer(players, matchName)
				if player == nil {
					player = &t.Player{
						Id:        strconv.Itoa(currId),
						Position:  pos,
						Name:      name,
						MatchName: matchName,
						FirstName: fName,
						LastName:  lName,
						Team:      team,
					}
					players = append(players, player)
					currId += 1
				}

				if rankType == "PPR" || rankType == "BOTH" {
					player.CustomPprRank = rank
				}
				if rankType == "STD" || rankType == "BOTH" {
					player.CustomStdRank = rank
				}
			}
		}
	})

	c.Visit(url)
	return players, currId
}
