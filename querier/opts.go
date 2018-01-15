package main

import (
	"flag"
)

func initFlags() {
	flag.StringVar(&outputFormat, "output", "basic",
		`desired output format. pick: basic, influx`)
	flag.BoolVar(&outputWithTimestamps, "timestamps", false,
		`include timestamps in influx line protocol output`)
	flag.Parse()
}
