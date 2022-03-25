package whois

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, err := json.Marshal(body)
	if err == nil {
		resp.Body = string(stringBody)
	} else {
		resp.Body = "Interface data could not be parsed."
	}
	return &resp, nil
}
