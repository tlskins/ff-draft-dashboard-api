package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

func SuccessResp(data interface{}, allowedOrigin string) (out Response, err error) {
	var buf bytes.Buffer

	body, err := json.Marshal(data)
	if err != nil {
		return Response{StatusCode: http.StatusUnprocessableEntity}, err
	}
	json.HTMLEscape(&buf, body)

	out = Response{
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

	return
}
