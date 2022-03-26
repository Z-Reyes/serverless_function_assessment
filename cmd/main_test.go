package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

//TestHandler tests handler method using valid and invalid input.
//Also acts as a test for whois.UnhandledMethod
func TestHandler(t *testing.T) {
	testCases := map[string]int{"GET": http.StatusOK, "POST": http.StatusMethodNotAllowed}

	var dummyRequest events.APIGatewayProxyRequest
	dummyRequest.RequestContext.Identity.SourceIP = "72.206.34.219"

	for key, val := range testCases {
		dummyRequest.HTTPMethod = key
		dummyResp, _ := handler(dummyRequest)
		if dummyResp.StatusCode != val {
			t.Fatal("Value of handler status code:", dummyResp.StatusCode, "does not match ground truth value:", val)
		}
	}
}
