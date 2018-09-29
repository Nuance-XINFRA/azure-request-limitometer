package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Nuance-Mobility/azure-request-limitometer/pkg/common"
	"github.com/Nuance-Mobility/azure-request-limitometer/pkg/outputs"

	"github.com/golang/glog"
)

const (
	cliName        = "limitometer"
	cliDescription = "Collects the number of remaining requests in Azure Resource Manager"
	cliVersion     = "2.0.0"
)

var config = common.Conf
var azureClient = common.Client

// flag variables
var nodename string
var outputTarget string
var useManagedIdentity bool
var flags map[string]*string

func init() {
	flag.StringVar(
		&outputTarget,
		"output",
		"influxdb",
		"Target output for the limitometer.")
	flag.StringVar(
		&nodename,
		"node",
		"",
		"Valid node in the resource group to create compute queries. Environment Variable: NODE_NAME")
	flag.Parse()
}

func main() {
	if len(flag.Args()) > 0 {
		if flag.Args()[0] == "help" {
			fmt.Printf("%s\n\n", cliName)
			fmt.Println(cliDescription)
			flag.PrintDefaults()
			os.Exit(2)
		}
		if flag.Args()[0] == "version" {
			fmt.Printf("%s version %s\n", cliName, cliVersion)
			os.Exit(0)
		}
	}
	if strings.ToLower(outputTarget) != "influxdb" {
		glog.Fatalln("That output is not yet implemented.")
		flag.PrintDefaults()
	}

	defer glog.Flush()

	if len(nodename) == 0 {
		nodename = os.Getenv("NODE_NAME")
		if len(nodename) == 0 {
			glog.Exit("Did not provide a nodename either through -node flag or Env Var NODE_NAME. Exiting.")
		}
	}

	glog.Infof("Starting limitometer with %s as target VM", nodename)

	requestsRemaining := getRequestsRemaining(nodename)

	glog.Infof("Writing to database: %s", outputTarget)
	outputs.WriteOutputInflux(requestsRemaining, "requestRemaining")
}
