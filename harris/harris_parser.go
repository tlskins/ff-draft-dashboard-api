package harris

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Position string

const (
	QB Position = "QB"
	RB Position = "RB"
	WR Position = "WR"
	TE Position = "TE"
)

type Player struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Name      string   `json:"name"`
	MatchName string   `json:"matchName"`
	Position  Position `json:"position"`
	Team      string   `json:"team"`

	ESPNRank      int `json:"espnRank"`
	HarrisRank    int `json:"harrisRank"`
	ESPNPPRRank   int `json:"espnPPRRank"`
	HarrisPPRRank int `json:"harrisPPRRank"`
}

func CleanName(name string) (out string) {
	cleanNmRgx := regexp.MustCompile(`\(.*\)`)
	out = strings.TrimSpace(cleanNmRgx.ReplaceAllString(name, ""))
	return
}

func MatchName(name string) (out string) {
	out = strings.ToLower(name)
	cleanRgx := regexp.MustCompile(`[^a-z ]`)
	out = cleanRgx.ReplaceAllString(out, "")
	out = strings.TrimSpace(out)
	return
}

func FindPlayer(players []*Player, matchName string) (out *Player) {
	matchParts := strings.Split(matchName, " ")
	matchFirst := matchParts[0]
	matchLast := strings.Join(matchParts[1:], "")
	for _, player := range players {
		nameParts := strings.Split(player.MatchName, " ")
		nameFirst := nameParts[0]
		nameLast := strings.Join(nameParts[1:], "")
		if (len(matchFirst) <= len(nameFirst) && !regexp.MustCompile(matchFirst).MatchString(nameFirst)) ||
			(len(nameFirst) < len(matchFirst) && !regexp.MustCompile(nameFirst).MatchString(matchFirst)) {
			continue
		}
		if (len(matchLast) <= len(nameLast) && !regexp.MustCompile(matchLast).MatchString(nameLast)) ||
			(len(nameLast) < len(matchLast) && !regexp.MustCompile(nameLast).MatchString(matchLast)) {
			continue
		}
		return player
	}
	return
}

func ParseHarrisRanks(url string, pos Position) (players []*Player) {
	players = []*Player{}
	isRankRgx := regexp.MustCompile(`^[0-9]+$`)
	isTmRgx := regexp.MustCompile(`^[A-Z]{2,}$`)
	isStdScrRgx := regexp.MustCompile(`(?i)standard scoring`)
	isPprScrRgx := regexp.MustCompile(`(?i)ppr scoring`)

	isPPR := false
	isCreate := true

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
				isPPR = false
				if len(players) != 0 {
					isCreate = false
				}
				continue
			}
			if isPprScrRgx.MatchString(text) {
				isPPR = true
				if len(players) != 0 {
					isCreate = false
				}
				continue
			}

			if !isTmRgx.MatchString(text) {
				name = CleanName(text)
				matchName = MatchName(text)
				nameParts := strings.Split(name, " ")
				fName = nameParts[0]
				lName = strings.Join(nameParts[1:], " ")
				continue
			} else {
				team = strings.TrimSpace(text)

				var player *Player
				if !isCreate {
					player = FindPlayer(players, matchName)
				}
				if player == nil {
					player = &Player{
						Position:  pos,
						Name:      name,
						MatchName: matchName,
						FirstName: fName,
						LastName:  lName,
						Team:      team,
					}
					players = append(players, player)
				}

				if isPPR {
					player.HarrisPPRRank = rank
				} else {
					player.HarrisRank = rank
				}
			}
		}
	})

	c.Visit(url)
	return
}

// func main() {
// 	qbs := ParseHarrisRanks("https://www.harrisfootball.com/ranks-draft", QB)
// 	wrs := ParseHarrisRanks("https://www.harrisfootball.com/wr-ranks-draft", WR)
// 	rbs := ParseHarrisRanks("https://www.harrisfootball.com/rb-ranks-draft", RB)
// 	tes := ParseHarrisRanks("https://www.harrisfootball.com/te-ranks-draft", TE)

// 	spew.Dump(qbs)
// 	spew.Dump(wrs)
// 	spew.Dump(rbs)
// 	spew.Dump(tes)
// }
