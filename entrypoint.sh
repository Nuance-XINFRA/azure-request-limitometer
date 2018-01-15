#!/bin/bash

set -e

./querier --output=influx "$1" > data
cat data
./pusher.sh influxdb:8086 k8s < data
