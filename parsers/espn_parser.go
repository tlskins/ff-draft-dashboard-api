package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	t "github.com/my_projects/ff-draft-dashboard-api/types"
	"github.com/pkg/errors"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:    60 * time.Second,
			DisableCompression: true,
		},
		Timeout: 60 * time.Second,
	}
}

func HttpRequest(client *http.Client, method, url string, headers map[string][]string, data, out interface{}) error {
	var body io.Reader
	if data != nil {
		if b, ok := data.([]byte); ok {
			body = bytes.NewReader(b)
		} else {
			if b, err := json.Marshal(data); err == nil {
				body = bytes.NewReader(b)
			} else {
				return errors.New(fmt.Sprintf("Error marshaling request data: %v", err.Error()))
			}
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error building http request: %v", err.Error()))
	}
	for header, values := range headers {
		req.Header[header] = values
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending http request: %v", err.Error()))
	}

	if out != nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading response body: %v", err.Error()))
		}
		if len(b) > 0 {
			// fmt.Println(string(b))
			return json.Unmarshal(b, out)
		}
	}

	return nil
}

func GetEspnApiUrl(year int) string {
	return fmt.Sprintf("https://fantasy.espn.com/apis/v3/games/ffl/seasons/%v/segments/0/leaguedefaults/3?view=kona_player_info", year)
}

func EspnQueryHeader(year, limit, offset int) (out map[string][]string) {
	out = make(map[string][]string)
	// hdr := fmt.Sprintf(`{"players":{"filterSlotIds":{"value":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,23,24]},"filterStatsForExternalIds":{"value":[2022,2023]},"filterStatsForSourceIds":{"value":[0,1]},"sortAdp":{"sortPriority":2,"sortAsc":true},"sortDraftRanks":{"sortPriority":2,"sortAsc":true,"value":"PPR"},"sortPercOwned":{"sortAsc":false,"sortPriority":4},"limit":%v,"offset":%v,"filterRanksForScoringPeriodIds":{"value":[1]},"filterRanksForRankTypes":{"value":["PPR"]},"filterRanksForSlotIds":{"value":[0,2,4,6,17,16]},"filterStatsForTopScoringPeriodIds":{"value":2,"additionalValue":["002023","102023","002022","022023"]}}}`, limit, offset)
	// hdr := fmt.Sprintf(`{"players":{"filterSlotIds":{"value":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,23,24]},"sortAdp":{"sortPriority":2,"sortAsc":true},"sortDraftRanks":{"sortPriority":100,"sortAsc":true,"value":"STANDARD"},"limit":%v,"offset":%v,"filterRanksForSlotIds":{"value":[0,2,4,6,17,16]},"filterStatsForTopScoringPeriodIds":{"value":2,"additionalValue":["002023","102023","002022","022023"]}}}`, limit, offset)
	hdr := fmt.Sprintf(`{"players":{"filterStatsForExternalIds":{"value":[%v,%v]},"filterSlotIds":{"value":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,23,24]},"filterStatsForSourceIds":{"value":[0,1]},"useFullProjectionTable":{"value":true},"sortAppliedStatTotal":{"sortAsc":false,"sortPriority":3,"value":"102023"},"sortDraftRanks":{"sortPriority":2,"sortAsc":true,"value":"PPR"},"sortPercOwned":{"sortPriority":4,"sortAsc":false},"limit":%v,"offset":%v,"filterRanksForSlotIds":{"value":[0,2,4,6,17,16]},"filterStatsForTopScoringPeriodIds":{"value":2,"additionalValue":["002023","102023","002022","022023"]}}}`, year-1, year, limit, offset)
	out["x-fantasy-filter"] = []string{hdr}
	return
}

func GetEspnPlayersForYear(client *http.Client, year int) ([]*t.EspnPlayer, error) {
	out := t.EspnPlayersResp{}
	// todo - add pagination
	if err := HttpRequest(client, "GET", GetEspnApiUrl(year), EspnQueryHeader(year, 350, 0), nil, &out); err != nil {
		return nil, err
	}

	return out.Players, nil
}
