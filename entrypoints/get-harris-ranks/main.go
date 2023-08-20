package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	client := p.NewHttpClient()

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
