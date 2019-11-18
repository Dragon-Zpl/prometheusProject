#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./alertmanager main.go
scp  -P 58422 ./alertmanager   root@192.168.188.34:/www/alertmanager/