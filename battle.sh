#!/usr/bin/env bash
curl -X POST -o match.json http://localhost:3008/create
./hal9k &
sleep 1
./hal9k &
curl -H "Content-Type: application/json" -d @match.json http://localhost:3008/start
trap "kill 0" SIGINT SIGTERM EXIT
wait
