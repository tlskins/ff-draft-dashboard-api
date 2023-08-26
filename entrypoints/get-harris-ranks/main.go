package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/my_projects/ff-draft-dashboard-api/api"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	s "github.com/my_projects/ff-draft-dashboard-api/store"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func Handler(ctx context.Context) (api.Response, error) {
	mongoDbName := os.Getenv("DRAFT_DASHBOARD_DB_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	mongoHost := os.Getenv("MONGO_HOST")

	client := p.NewHttpClient()
	store, err := s.NewStore(mongoDbName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		return api.Response{StatusCode: http.StatusInternalServerError}, err
	}
	reportsCol := store.PlayerReportsCol()
	now := time.Now()
	var year int
	if year, err = strconv.Atoi(now.Format("2006")); err != nil {
		return api.Response{StatusCode: http.StatusInternalServerError}, err
	}

	// get all espn players
	players, err := p.GetEspnPlayersForYear(client, year, 350)
	if err != nil {
		return api.Response{StatusCode: http.StatusInternalServerError}, err
	}
	fmt.Printf("found %v players\n", len(players))

	// get custom ranks
	harrisPlayers := p.ParseHarrisRanksV2(year)
	fmt.Printf("found %v harris players\n", len(harrisPlayers))
	p.AddHarrisRanks(harrisPlayers, players)

	// fetch all player reports and add to a map
	playerReports := []*t.PlayerReport{}
	if err = store.Find(reportsCol, s.M{}, &playerReports); err != nil {
		return api.Response{StatusCode: http.StatusInternalServerError}, err
	}
	playerReportsMap := map[string]*t.PlayerReport{}
	for _, report := range playerReports {
		playerReportsMap[report.Id] = report
	}

	// add player report data
	for _, player := range players {
		report := playerReportsMap[player.Id]
		if report != nil {
			player.AddPlayerReport(report)
		}
	}

	// calc stats by num teams
	posStatsByNumTeamByYear := p.CalcAllStats(players, year)

	resp, err := api.SuccessResp(
		map[string]interface{}{"players": players, "posStatsByNumTeamByYear": posStatsByNumTeamByYear},
		os.Getenv("ALLOWED_ORIGIN"),
	)

	return resp, err
}

func main() {
	lambda.Start(Handler)
}
