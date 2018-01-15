#!/bin/bash

set -e

# expects input from az-limitometer/querier
[[ "$#" == 2 ]] || ( echo "usage: pusher <server> <database>" && exit 1 )

server="$1"
database="$2"
data="$(cat /dev/stdin)"

curl -XPOST "http://$server/write?db=$database" --data-binary "$data"
