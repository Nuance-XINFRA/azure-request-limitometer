package main

import (
	"fmt"
	"time"
)

// looks like:
// [root@kn-default-0 ~]# HTTP_PROXY=10.33.59.5:3128 ./querier `hostname -s`
// Microsoft.Compute/HighCostGet30Min 1365
// Microsoft.Compute/HighCostGet3Min 272
// Microsoft.Compute/LowCostGet3Min 3398
// Microsoft.Compute/LowCostGet30Min 23785
func writeOutputBasic(requestsRemaining map[string]int) {
	for k, v := range requestsRemaining {
		fmt.Println(k, v)
	}
}

// looks like:
// [root@kn-default-0 ~]# HTTP_PROXY=10.33.59.5:3128 ./querier --output-influx `hostname -s`
// Microsoft.Compute/LowCostGet3Min requestsRemaining=3244 1513882044
// Microsoft.Compute/LowCostGet30Min requestsRemaining=23353 1513882044
// Microsoft.Compute/HighCostGet3Min requestsRemaining=271 1513882044
// Microsoft.Compute/HighCostGet30Min requestsRemaining=1364 1513882044
func writeOutputInflux(requestsRemaining map[string]int) {
	if outputWithTimestamps {
		// TODO maybe this timestamp should be taken before the request
		nowTimestamp := time.Now().Unix()
		for k, v := range requestsRemaining {
			fmt.Printf("%v requestsRemaining=%v %v\n", k, v, nowTimestamp)
		}
	} else {
		for k, v := range requestsRemaining {
			fmt.Printf("%v requestsRemaining=%v\n", k, v)
		}
	}
}
