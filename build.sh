#! /bin/bash

go build -o dist/main
server_file=server.js
start_go_cmd="DEBUG=t PROJECTROOT=/Users/Natalzia/go_work/src/github.com/jnatalzia/gralhund-bot ./dist/main -t $AUTHCODE >> /tmp/log.log &"
echo $start_go_cmd
 
./kill_gralhund.sh
 
echo "building binary"
go build -o dist/main
echo "starting server"
eval "$start_go_cmd"  