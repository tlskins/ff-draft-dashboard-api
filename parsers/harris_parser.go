package parsers

import (
	"fmt"
	"log"
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

func ParseHarrisName(name string) (fName, lName string) {
	nameParts := strings.Split(name, " ")
	fName = nameParts[0]
	lName = strings.Join(nameParts[1:], " ")

	return
}

func ParseHarrisTeam(team string) (out string) {
	switch team {
	case "PHI":
		return "PHL"
	default:
		return team
	}
}

func HarrisPlayerPrimaryKey(player *t.HarrisPlayer) string {
	alphaRgx := regexp.MustCompile("[^a-zA-Z]+")
	nameKey := strings.ToUpper(alphaRgx.ReplaceAllString(player.Name, ""))
	return fmt.Sprintf("%s-%s-%s", nameKey, strings.ToUpper(player.Team), strings.ToUpper(player.Position))
}

func ParsePlayersForPath(pos, urlPath string, isRankByType bool) (players []*t.HarrisPlayer) {
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
							newPlayer.FirstName, newPlayer.LastName = ParseHarrisName(newPlayer.Name)
						} else if idx == 2 {
							// team
							newPlayer.Team = ParseHarrisTeam(cellText)
						}
					}
				})

				// player parsing complete
				// fmt.Sprintf("found %v, ppr %v, rank %v\n", playerFound, isPpr, rank)
				if playerFound && isPpr != nil && rank > 0 {
					newPlayer.Id = HarrisPlayerPrimaryKey(&newPlayer)

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

func ParseHarrisRanksV2(year int) (out []*t.HarrisPlayer) {
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
	out = []*t.HarrisPlayer{}
	for pos, urlPath := range posToPathMap {
		posPlayers := ParsePlayersForPath(pos, urlPath, posToRankTypeMap[pos])
		out = append(out, posPlayers...)
	}

	return
}

// mutates input players to include harris ranks
func AddHarrisRanks(harrisPlayers []*t.HarrisPlayer, players []*t.Player) (unmatched []*t.Player) {
	matches := []*t.HarrisPlayerMatch{}
	unmatched = []*t.Player{}

	// build harris lookups
	harrisMap := map[string]*t.HarrisPlayer{}
	harrisMatched := map[string]bool{}
	harrisTeamPosLookup := map[string]map[string][]*t.HarrisPlayer{}
	for _, harrisPlayer := range harrisPlayers {
		harrisMap[harrisPlayer.Id] = harrisPlayer
		team := harrisPlayer.Team
		pos := harrisPlayer.Position
		if harrisTeamPosLookup[team] == nil {
			harrisTeamPosLookup[team] = make(map[string][]*t.HarrisPlayer)
		}
		if harrisTeamPosLookup[team][pos] == nil {
			harrisTeamPosLookup[team][pos] = []*t.HarrisPlayer{}
		}
		harrisTeamPosLookup[team][pos] = append(harrisTeamPosLookup[team][pos], harrisPlayer)
	}

	// match espn to harris players
	for _, player := range players {
		if player.Position == t.DST || player.Position == t.NoPosition {
			continue
		}
		playerMatch := &t.HarrisPlayerMatch{Player: player}
		matches = append(matches, playerMatch)
		team := player.Team
		pos := string(player.Position)

		var matchedHarrisPlayer *t.HarrisPlayer
		// primary match via direct user lookup
		playerMatchKey := player.MatchKey()
		fmt.Printf("harriskey: %s\n", playerMatchKey)
		if harrisMap[playerMatchKey] != nil {
			harrisMatched[playerMatchKey] = true
			matchedHarrisPlayer = harrisMap[playerMatchKey]
		} else {
			// secondary match by Levenshtein distance algorithm
			isSecondaryMatch := false
			for _, harrisPlayer := range harrisTeamPosLookup[team][pos] {
				if harrisMatched[harrisPlayer.Id] {
					continue
				}
				diffScore := StringDiffScore(playerMatchKey, harrisPlayer.Id)
				log.Printf("\tdiff score: %v %s\n", diffScore, harrisPlayer.Name)
				if diffScore <= 5 {
					isSecondaryMatch = true
					log.Printf("\tMatching: %s and %s with\n", player.Name, harrisPlayer.Name)
					matchedHarrisPlayer = harrisPlayer
				}
			}

			if !isSecondaryMatch {
				log.Printf("NOT MATCHED %s %s %v\n", playerMatchKey, player.Name, player.EspnOvrStdRank)
			}
		}
		if matchedHarrisPlayer != nil {
			playerMatch.Harris = matchedHarrisPlayer
		}
	}

	// mutate players and add ranks
	for _, match := range matches {
		match.AddPlayerRank()
	}

	return
}
