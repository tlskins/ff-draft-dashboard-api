package parsers

import (
	"fmt"
	"net/http"

	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func GetEspnApiUrl(year int) string {
	return fmt.Sprintf("https://fantasy.espn.com/apis/v3/games/ffl/seasons/%v/segments/0/leaguedefaults/3?view=kona_player_info", year)
}

func EspnQueryHeader(year, limit, offset int) (out map[string][]string) {
	out = make(map[string][]string)
	// hdr := fmt.Sprintf(`{"players":{"filterSlotIds":{"value":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,23,24]},"filterStatsForExternalIds":{"value":[2022,2023]},"filterStatsForSourceIds":{"value":[0,1]},"sortAdp":{"sortPriority":2,"sortAsc":true},"sortDraftRanks":{"sortPriority":2,"sortAsc":true,"value":"PPR"},"sortPercOwned":{"sortAsc":false,"sortPriority":4},"limit":%v,"offset":%v,"filterRanksForScoringPeriodIds":{"value":[1]},"filterRanksForRankTypes":{"value":["PPR"]},"filterRanksForSlotIds":{"value":[0,2,4,6,17,16]},"filterStatsForTopScoringPeriodIds":{"value":2,"additionalValue":["002023","102023","002022","022023"]}}}`, limit, offset)
	// hdr := fmt.Sprintf(`{"players":{"filterSlotIds":{"value":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,23,24]},"sortAdp":{"sortPriority":2,"sortAsc":true},"sortDraftRanks":{"sortPriority":100,"sortAsc":true,"value":"STANDARD"},"limit":%v,"offset":%v,"filterRanksForSlotIds":{"value":[0,2,4,6,17,16]},"filterStatsForTopScoringPeriodIds":{"value":2,"additionalValue":["002023","102023","002022","022023"]}}}`, limit, offset)
	hdr := fmt.Sprintf(`{"players":{"filterStatsForExternalIds":{"value":[%v,%v]},"filterSlotIds":{"value":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,23,24]},"filterStatsForSourceIds":{"value":[0,1]},"useFullProjectionTable":{"value":true},"sortAppliedStatTotal":{"sortAsc":false,"sortPriority":3,"value":"102023"},"sortDraftRanks":{"sortPriority":2,"sortAsc":true,"value":"PPR"},"sortPercOwned":{"sortPriority":4,"sortAsc":false},"limit":%v,"offset":%v,"filterRanksForSlotIds":{"value":[0,2,4,6,17,16]},"filterStatsForTopScoringPeriodIds":{"value":2,"additionalValue":["002023","102023","022023","002022"]}}}`, year-1, year, limit, offset)
	out["x-fantasy-filter"] = []string{hdr}
	return
}

func GetEspnPlayersForYear(client *http.Client, year, numPlayers int) (players []*t.Player, err error) {
	espnPlayersResp := t.EspnPlayersResp{}
	if err = HttpRequest(client, "GET", GetEspnApiUrl(year), EspnQueryHeader(year, numPlayers, 0), nil, &espnPlayersResp); err != nil {
		return nil, err
	}

	players = make([]*t.Player, len(espnPlayersResp.Players))
	for i, espnPlayer := range espnPlayersResp.Players {
		players[i] = espnPlayer.ToPlayer()
	}

	return
}
