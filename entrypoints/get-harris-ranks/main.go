package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	h "github.com/my_projects/ff-draft-dashboard-api/harris"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	qbs := h.ParseHarrisRanks("https://www.harrisfootball.com/ranks-draft", h.QB)
	wrs := h.ParseHarrisRanks("https://www.harrisfootball.com/wr-ranks-draft", h.WR)
	rbs := h.ParseHarrisRanks("https://www.harrisfootball.com/rb-ranks-draft", h.RB)
	tes := h.ParseHarrisRanks("https://www.harrisfootball.com/te-ranks-draft", h.TE)

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

	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

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
