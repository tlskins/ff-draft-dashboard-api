package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	client := p.NewHttpClient()

	// get espn ranks
	out := t.EspnPlayersResp{}
	if err := p.HttpRequest(client, "GET", p.EspnApiUrl, p.EspnQueryHeader(350, 0), nil, &out); err != nil {
		return Response{StatusCode: 404}, err
	}
	players := []*t.Player{}
	for i, p := range out.Players {
		player := p.ToPlayer()
		if player.Position == t.NoPosition {
			continue
		}
		player.EspnAdp = i + 1
		players = append(players, player)
	}

	// get fpros ranks
	fprosOut, err := p.HttpHtmlRequest(client, "GET", p.FProsApiUrl, map[string][]string{}, nil)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	rgx := regexp.MustCompile(`var ecrData = ({.*})`)
	rs := rgx.FindStringSubmatch(fprosOut)
	byt := []byte(rs[1])

	fprosResp := t.FproEcrData{}
	if err := json.Unmarshal(byt, &fprosResp); err != nil {
		return Response{StatusCode: 404}, err
	}

	for _, p := range fprosResp.Players {
		matchName := t.MatchName(p.PlayerName)
		player := t.FindPlayer(players, matchName)
		if player == nil {
			player = p.ToPlayer()
			players = append(players, player)
		} else {
			player.CustomPprRank = p.RankEcr
			player.CustomStdRank = p.RankEcr
			if player.Tier == "" {
				player.Tier = strconv.Itoa(p.Tier)
			}
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
