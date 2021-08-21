package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	client := p.NewHttpClient()
	out := t.EspnPlayersResp{}
	if err := p.HttpRequest(client, "GET", p.EspnApiUrl, p.EspnQueryHeader(250, 0), nil, &out); err != nil {
		return Response{StatusCode: 404}, err
	}
	players := make([]*t.Player, len(out.Players))
	for i, p := range out.Players {
		players[i] = p.ToPlayer()
		players[i].EspnAdp = i + 1
	}

	currId := 1
	players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/ranks-draft", t.QB, currId, players)
	players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/wr-ranks-draft", t.WR, currId, players)
	players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/rb-ranks-draft", t.RB, currId, players)
	players, currId = p.ParseHarrisRanks("https://www.harrisfootball.com/te-ranks-draft", t.TE, currId, players)

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
