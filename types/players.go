package types

import (
	"regexp"
	"strings"
)

type Position string

const (
	QB         Position = "QB"
	RB         Position = "RB"
	WR         Position = "WR"
	TE         Position = "TE"
	NoPosition Position = ""
)

type Player struct {
	Id        int      `json:"id"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Name      string   `json:"name"`
	MatchName string   `json:"matchName"`
	Position  Position `json:"position"`
	Team      string   `json:"team"`
	Tier      string   `json:"tier"`

	CustomStdRank  int `json:"customStdRank,omitempty"`
	CustomPprRank  int `json:"customPprRank,omitempty"`
	EspnOvrPprRank int `json:"espnOvrPprRank,omitempty"`
	EspnOvrStdRank int `json:"espnOvrStdRank,omitempty"`
	EspnAdp        int `json:"espnAdp,omitempty"`
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
