package main

import (
	"azure-request-limitometer/common"
	"flag"
	"log"
)

var config = common.Conf
var client = common.Client

// flag variables
var outputFormat string
var outputWithTimestamps bool

func init() {
	initFlags()
}

func main() {
	if len(flag.Args()) != 1 {
		// NOTE requests remaining does not depend on
		// the nodename used, but it must be valid
		log.Fatalln("usage: querier <nodename>")
	}

	nodename := flag.Args()[0]
	requestsRemaining := getRequestsRemaining(nodename)

	switch outputFormat {
	case "basic":
		writeOutputBasic(requestsRemaining)
	case "influx":
		writeOutputInflux(requestsRemaining)
	default:
		log.Fatalln("unknown output format: %v", outputFormat)
	}
}
