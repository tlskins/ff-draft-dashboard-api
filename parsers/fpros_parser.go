package parsers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func GetFprosUrl(isPpr bool) string {
	if isPpr {
		// https://www.fantasypros.com/nfl/rankings/ppr.php
		return "https://www.fantasypros.com/nfl/rankings/ppr-cheatsheets.php"
	} else {
		return "https://www.fantasypros.com/nfl/rankings/consensus-cheatsheets.php"
	}
}

func GetFProsPlayersByFormat(client *http.Client, isPpr bool) (out *t.FproEcrData, err error) {
	var data string
	data, err = HttpHtmlRequest(client, "GET", GetFprosUrl(isPpr), map[string][]string{}, nil)
	if err != nil {
		return
	}

	rgx := regexp.MustCompile(`var ecrData = ({.*});`)
	rs := rgx.FindStringSubmatch(data)
	byt := []byte(rs[1])

	out = &t.FproEcrData{}
	if err = json.Unmarshal(byt, out); err != nil {
		return
	}

	return
}

func GetFprosPlayers(client *http.Client) (out []*t.FproPlayer, err error) {
	out = []*t.FproPlayer{}
	playersDict := map[string]*t.FproPlayer{}

	pprPlayersData, err := GetFProsPlayersByFormat(client, true)
	if err != nil {
		return out, err
	}
	fmt.Printf("found %v ppr players\n", len(pprPlayersData.Players))
	for _, player := range pprPlayersData.Players {
		player.PprRank = player.RankEcr
		player.PosPprRank = player.PosRank
		playersDict[player.SportsDataId] = player
	}

	stdPlayersData, err := GetFProsPlayersByFormat(client, false)
	if err != nil {
		return out, err
	}
	fmt.Printf("found %v std players\n", len(stdPlayersData.Players))
	unmatched := 0
	for _, player := range stdPlayersData.Players {
		if playersDict[player.SportsDataId] == nil {
			playersDict[player.SportsDataId] = player
			unmatched += 1
			fmt.Println(player.PlayerName)
		}
		playersDict[player.SportsDataId].StdRank = player.RankEcr
		playersDict[player.SportsDataId].PosStdRank = player.PosRank
	}
	log.Printf("unmatched = %v\n", unmatched)

	for _, player := range playersDict {
		out = append(out, player)
	}

	return
}

func AddFprosRanks(fprosPlayers []*t.FproPlayer, players []*t.Player) (unmatched []*t.Player) {
	matches := []*t.FprosPlayerMatch{}
	unmatched = []*t.Player{}

	// build fpros lookups
	fprosMap := map[string]*t.FproPlayer{}
	fprosMatched := map[string]bool{}
	fprosTeamPosLookup := map[string]map[string][]*t.FproPlayer{}
	for _, fprosPlayer := range fprosPlayers {
		fprosMap[fprosPlayer.PrimaryKey()] = fprosPlayer
		team := fprosPlayer.PlayerTeamId
		pos := fprosPlayer.PlayerPositionId
		if fprosTeamPosLookup[team] == nil {
			fprosTeamPosLookup[team] = make(map[string][]*t.FproPlayer)
		}
		if fprosTeamPosLookup[team][pos] == nil {
			fprosTeamPosLookup[team][pos] = []*t.FproPlayer{}
		}
		fprosTeamPosLookup[team][pos] = append(fprosTeamPosLookup[team][pos], fprosPlayer)
	}

	// match espn to fpros players
	for _, player := range players {
		if player.Position == t.DST || player.Position == t.NoPosition {
			continue
		}
		playerMatch := &t.FprosPlayerMatch{Player: player}
		matches = append(matches, playerMatch)
		team := player.Team
		pos := string(player.Position)

		var matchedFproPlayer *t.FproPlayer
		// primary match via direct user lookup
		playerMatchKey := player.MatchKey()
		if fprosMap[playerMatchKey] != nil {
			fprosMatched[playerMatchKey] = true
			matchedFproPlayer = fprosMap[playerMatchKey]
		} else {
			// secondary match by Levenshtein distance algorithm
			isSecondaryMatch := false
			for _, fprosPlayer := range fprosTeamPosLookup[team][pos] {
				if fprosMatched[fprosPlayer.SportsDataId] {
					continue
				}
				diffScore := StringDiffScore(playerMatchKey, fprosPlayer.SportsDataId)
				log.Printf("\tdiff score: %v %s\n", diffScore, fprosPlayer.PlayerName)
				if diffScore <= 5 {
					isSecondaryMatch = true
					log.Printf("\tMatching: %s and %s with\n", player.Name, fprosPlayer.PlayerName)
					matchedFproPlayer = fprosPlayer
				}
			}

			if !isSecondaryMatch {
				log.Printf("NOT MATCHED %s %s %v\n", playerMatchKey, player.Name, player.EspnOvrStdRank)
			}
		}
		if matchedFproPlayer != nil {
			playerMatch.Fpros = matchedFproPlayer
		}
	}

	// mutate players and add ranks
	for _, match := range matches {
		match.AddPlayerRank()
	}

	return
}
