package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/golang/glog"
)

// Example Request Headers:
// 'x-ms-ratelimit-remaining-resource': 'Microsoft.Compute/HighCostGet3Min;133,Microsoft.Compute/HighCostGet30Min;657'
// 'x-ms-ratelimit-remaining-resource': 'Microsoft.Compute/LowCostGet3Min;3989,Microsoft.Compute/LowCostGet30Min;31790'
// 'x-ms-ratelimit-remaining-resource': 'Microsoft.Compute/PutVM3Min;740,Microsoft.Compute/PutVM30Min;3695'

var expectedHeaderField = "X-Ms-Ratelimit-Remaining-Resource"
var expectedHeaderFormat = regexp.MustCompile(`(Microsoft.Compute/\w+);(\d+)`)

func getRequestsRemaining(nodename string) (requestsRemaining map[string]int) {
	requestsRemaining = make(map[string]int)

	responses := []autorest.Response{
		azureClient.GetVM(nodename).Response,
		azureClient.GetAllVM().Response().Response,
		azureClient.PutVM(nodename),
	}

	for _, response := range responses {
		if response.StatusCode != 200 {
			glog.Fatalf("Response did not return a StatusCode of 200. Check HTTP_PROXY. StatusCode: %d", response.StatusCode)
		}
		for k, v := range extractRequestsRemaining(response.Header) {
			requestsRemaining[k] = v
		}
	}

	return
}

func extractRequestsRemaining(h http.Header) (requestsRemaining map[string]int) {
	requestsRemaining = map[string]int{}

	headerSubfields := strings.Split(h.Get(expectedHeaderField), ",")

	for _, field := range headerSubfields {
		matches := expectedHeaderFormat.FindStringSubmatch(field)
		if !(len(matches) == 3) {
			glog.Errorf("header didn't contain expected data: %s", field)
		}

		requestType := matches[1]
		requestsLeft, err := strconv.Atoi(matches[2])
		if err != nil {
			glog.Error(err)
		}
		requestsRemaining[requestType] = requestsLeft
	}

	return requestsRemaining
}
