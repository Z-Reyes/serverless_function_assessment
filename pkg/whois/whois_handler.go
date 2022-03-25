package whois

import (
	"fmt"
	"net"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/likexian/whois"
)

const (
	//path parameter key (see main.tf)
	requestedIP = "requestip"
	queryParam  = "trim"

	//generic error message string
	errorMethodNotAllowed = "Method Not Allowed"

	regexOrg   = "OrgName:.*?\n"
	regexRange = "NetRange:.*?\n"
	regexAddr  = "Address:.*?\nCity:.*?\nStateProv:.*?\nPostalCode:.*?\nCountry:.*?\n"
)

type trimmedWhoIs struct {
	Org       string `json:"Organization"`
	AddrRange string `json:"Network Range"`
	Addr      string `json:"Address"`
}
type errorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

//Default response for when receiving unimplemented method
func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, errorMethodNotAllowed)
}

//Ip attempts to retrieve whoIS IP data based on an incoming APIGateway Proxy Request.
//If no path parameter is included, it provides info about the source IP that sent the request.
//If a path parameter is included and it's a valid IP address (ipv4/6), it attempts
//to provide whoIS data about that IP address.
//A query parameter of '?trim=true' indicates that we need to provide a truncated response.
func Ip(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse,
	error) {

	//Determine if path parameter exists. For this system, the path parameter is 'requestid' (see main.tf)
	var searchTerm string
	if val, ok := req.PathParameters[requestedIP]; ok {
		searchTerm = val
	} else {
		searchTerm = req.RequestContext.Identity.SourceIP
	}

	//Check search term for valid address
	if !isValidIpAddress(searchTerm) {
		return apiResponse(http.StatusBadRequest, errorBody{aws.String(fmt.Sprintf("Invalid path parameter: %s. Please use a valid PUBLIC IPv4 or IPv6 address", searchTerm))})
	}
	result, err := whois.Whois(searchTerm)
	if err != nil {
		return apiResponse(http.StatusBadRequest, errorBody{aws.String(err.Error())})
	}

	//Check if query parameter exists. If it does, trim data.
	if val, ok := req.QueryStringParameters[queryParam]; ok {
		if len(val) > 0 && val == "true" {
			return apiResponse(http.StatusOK, trimWhoIs(result))
		}
	}
	return apiResponse(http.StatusOK, result)
}

func isValidIpAddress(ip string) bool {

	isIP := net.ParseIP(ip)
	if isIP == nil {
		return false
	}
	return !(isIP.IsPrivate())
}

func trimWhoIs(input string) trimmedWhoIs {
	//Retrieve 3 pieces of information: Organization, Netrange, address
	var trimmed trimmedWhoIs
	r, err := regexp.Compile(regexOrg)
	if err != nil {
		trimmed.Org = "REGEX ERR: Organization expression could not be compiled."
	}
	trimmed.Org = r.FindString(input)
	r, err = regexp.Compile(regexAddr)
	if err != nil {
		trimmed.Addr = "REGEX ERR: Address expression could not be compiled."
	}
	trimmed.Addr = r.FindString(input)
	r, err = regexp.Compile(regexRange)
	if err != nil {
		trimmed.AddrRange = "REGEX ERR: Network Range expression could not be compiled."
	}
	trimmed.AddrRange = r.FindString(input)
	if trimmed.Org == "" && trimmed.Addr == "" && trimmed.AddrRange == "" {
		trimmed.Org = "OrgName token not found."
		trimmed.Addr = "Address token not found."
		trimmed.AddrRange = "Network Address token not found."
	}
	return trimmed
}
