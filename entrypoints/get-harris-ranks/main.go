package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	h "github.com/my_projects/ff-draft-dashboard-api/harris"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	currId := 1
	var qbs, wrs, rbs, tes []*h.Player
	qbs, currId = h.ParseHarrisRanks("https://www.harrisfootball.com/ranks-draft", h.QB, currId, false)
	wrs, currId = h.ParseHarrisRanks("https://www.harrisfootball.com/wr-ranks-draft", h.WR, currId, false)
	rbs, currId = h.ParseHarrisRanks("https://www.harrisfootball.com/rb-ranks-draft", h.RB, currId, false)
	tes, currId = h.ParseHarrisRanks("https://www.harrisfootball.com/te-ranks-draft", h.TE, currId, false)

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		string(h.QB): qbs,
		string(h.WR): wrs,
		string(h.RB): rbs,
		string(h.TE): tes,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      allowedOrigin,
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
