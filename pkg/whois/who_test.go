package whois

import (
	"math"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

//TestAPIResponse tests APIResponse function on valid and message bodies
func TestAPIResponse(t *testing.T) {
	validTestData := "String data"
	invalidTestDataType := make(chan int)
	invalidTestDataValue := math.Inf(1)

	_, err := apiResponse(1, validTestData)
	if err != nil {
		t.Fatal("Valid data returned error")
	}
	_, err = apiResponse(1, invalidTestDataType)
	if err == nil {
		t.Fatal("Invalid data type did NOT return an error")
	}
	_, err = apiResponse(1, invalidTestDataValue)
	if err == nil {
		t.Fatal("Invalid data value did NOT return an error")
	}
}

//TestIsValidIpAddress tests isValidIpAddress using valid and invalid input
func TestIsValidIpAddress(t *testing.T) {
	result := map[string]bool{"abcdefg": false, "72.206.34.219": true,
		"8.8.8.8": true, "182.1.1.2": true,
		"222.222.222.222": true, "192.168.1.1": false,
		"172.16.1.1": false, "10.1.1.1": false,
		"2001:4860:4860::8888": true, "fc00::": false}

	for key, val := range result {
		if val != isValidIpAddress(key) {
			t.Fatal("Value of isValidIpAddress:", isValidIpAddress(key), "for key", key, ",does not match ground truth value:", val)
		}
	}
}

//TestTrimWhoIs tests trimWhoIs using valid and invalid input
func TestTrimWhoIs(t *testing.T) {
	testString := "aaabbbccc"
	orgResults := map[string]string{"aaa": "aaa", "cab": "OrgName token not found.", `^\/(?!\/)(.*?)`: "REGEX ERR: Organization expression could not be compiled."}
	addrResults := map[string]string{"bbb": "bbb", "cab": "Address token not found.", `^\/(?!\/)(.*?)`: "REGEX ERR: Address expression could not be compiled."}
	rangeResults := map[string]string{"ccc": "ccc", "cab": "Network Address token not found.", `^\/(?!\/)(.*?)`: "REGEX ERR: Network Address expression could not be compiled."}
	for key, val := range orgResults {
		testValue := trimWhoIs(testString, key, "", "")
		if val != testValue.Org {
			t.Fatal("Value of trimWhoIs.Org:", testValue.Org, "for key", key, ",does not match ground truth value:", val)
		}
	}
	for key, val := range addrResults {
		testValue := trimWhoIs(testString, "", key, "")
		if val != testValue.Addr {
			t.Fatal("Value of trimWhoIs.Addr:", testValue.Addr, "for key", key, ",does not match ground truth value:", val)
		}
	}
	for key, val := range rangeResults {
		testValue := trimWhoIs(testString, "", "", key)
		if val != testValue.AddrRange {
			t.Fatal("Value of trimWhoIs.AddrRange:", testValue.AddrRange, "for key", key, ",does not match ground truth value:", val)
		}
	}
}

type ipSelection struct {
	status      int
	isPathSpoof bool
}

//TestIp tests Ip function using valid and invalid input
func TestIp(t *testing.T) {
	dummyRequest := events.APIGatewayProxyRequest{}

	reqs := map[string]ipSelection{"72.206.34.219": ipSelection{http.StatusOK, false}, "72.206.34.220": ipSelection{http.StatusOK, true},
		"192.168.1.17": ipSelection{http.StatusBadRequest, false}, "192.168.1.18": ipSelection{http.StatusBadRequest, true},
		"2001:4860:4860::8888": ipSelection{http.StatusOK, false}, "2001:4860:4860::8889": ipSelection{http.StatusOK, true},
		"fc00::": ipSelection{http.StatusBadRequest, false}, "fbff::": ipSelection{http.StatusBadRequest, true}}

	for key, val := range reqs {
		dummyRequest.RequestContext.Identity.SourceIP = ""
		dummyRequest.PathParameters = make(map[string]string)
		if val.isPathSpoof {
			dummyRequest.PathParameters[requestedIP] = key
		} else {
			dummyRequest.RequestContext.Identity.SourceIP = key
		}
		res, _ := Ip(dummyRequest)
		if res.StatusCode != val.status {
			t.Fatal("Value of Ip status code:", res.StatusCode, "does not match ground truth value:", val.status)
		}
	}
}
