package main

import (
	"serverless_function_golang/pkg/whoishandler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	/*region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)
	*/
	lambda.Start(handler)
}

//const tableName = "LambdaInGoUser"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return whoishandler.GetWhoIs(req)
		/*
			case "POST":
				return handlers.CreateUser(req, tableName, dynaClient)
			case "PUT":
				return handlers.UpdateUser(req, tableName, dynaClient)
			case "DELETE":
				return handlers.DeleteUser(req, tableName, dynaClient)
		*/
	default:
		return whoishandler.UnhandledMethod()
	}
}
