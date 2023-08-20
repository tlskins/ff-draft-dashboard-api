package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/my_projects/ff-draft-dashboard-api/api"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	"github.com/my_projects/ff-draft-dashboard-api/store"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	mongoDbName := os.Getenv("DRAFT_DASHBOARD_DB_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	mongoHost := os.Getenv("MONGO_HOST")

	client := p.NewHttpClient()
	store, err := store.NewStore(mongoDbName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError}, err
	}
	playerReportsCol := store.PlayerReportsCol()

	espnPlayers, err := p.GetEspnPlayersForYear(client, 2023)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError}, err
	}
	fmt.Printf("found %v espn players\n", len(espnPlayers))
	harrisPlayers := p.ParseHarrisRanksV2(2023)
	fmt.Printf("found %v harris players\n", len(harrisPlayers))

	players, _, err := p.MatchHarrisAndEspnPlayers(harrisPlayers, espnPlayers)
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	// add player report data
	for _, player := range players {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result := playerReportsCol.FindOne(ctx, api.M{"_id": player.Id})
		if result.Err() == nil {
			report := &t.PlayerReport{}
			if err = result.Decode(report); err != nil {
				return Response{StatusCode: http.StatusInternalServerError}, err
			}
			player.Pros = report.Pros
			player.Cons = report.Cons
		}
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{"players": players})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      os.Getenv("ALLOWED_ORIGIN"),
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "OPTIONS,POST,GET",
			"Access-Control-Allow-Headers":     "Content-Type",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
