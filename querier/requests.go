package main

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func getRequestsRemaining(nodename string) (requestsRemaining map[string]int) {
	requestsRemaining = map[string]int{}
	reponses := []autorest.Response{
		client.GetVM(nodename).Response, // low cost get
		client.GetAllVM().Response,      // high cost get
		client.PutVM(nodename).Response, //  put vm
		// TODO disk operations
	}

	for _, response := range reponses {
		assert(response.StatusCode < 300, fmt.Sprintf("VMClient GET failed: %v", response.Status))

		for k, v := range extractRequestsRemaining(response.Header) {
			requestsRemaining[k] = v
		}
	}

	return
}

// hard-coded
var expectedHeaderField = "X-Ms-Ratelimit-Remaining-Resource"
var expectedHeaderFormat = regexp.MustCompile(`(Microsoft.Compute/\w+);(\d+)`)

func extractRequestsRemaining(h http.Header) (requestsRemaining map[string]int) {
	requestsRemaining = map[string]int{}

	headerSubfields := strings.Split(h[expectedHeaderField][0], ",")

	for _, field := range headerSubfields {
		matches := expectedHeaderFormat.FindStringSubmatch(field)
		assert(len(matches) == 3, "header didn't contain expected data: "+field)

		requestType := matches[1]
		requestsLeft, err := strconv.Atoi(matches[2])
		xk(err)

		requestsRemaining[requestType] = requestsLeft
	}

	return requestsRemaining
}
