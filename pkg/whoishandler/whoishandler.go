package whoishandler

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/likexian/whois"
)

func GetWhoIs(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse,
	error) {
	var results []interface{}

	results = append(results, req)

	result, err := whois.Whois("72.206.34.219")
	if err == nil {
		//	fmt.Println(result)
		results = append(results, result)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	return apiResponse(http.StatusOK, results)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}

var ErrorMethodNotAllowed = "Method Not Allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}
