#!/usr/bin/env bash
./hal9k &
sleep 1
./hal9k &
curl -H "Content-Type: application/json" -d '{"match":"68267d65-6f63-4c18-8afa-dce5c91e0e73"}' http://localhost:3008/start
trap "kill 0" SIGINT SIGTERM EXIT
wait
