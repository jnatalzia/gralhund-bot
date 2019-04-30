#! /bin/bash

trap "./kill_gralhund.sh" INT

./build.sh

fswatch -o ./ -e ".*" -i "\\.go$" -i "\\.csv$" | xargs -n1 './build.sh'