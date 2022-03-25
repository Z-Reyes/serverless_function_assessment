//main contains our entry point for our serverless function.
//Simply connects lambda.Start to handler function.
package main

import (
	"serverless_function_golang/pkg/whois"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return whois.Ip(req)
	default:
		return whois.UnhandledMethod()
	}
}
