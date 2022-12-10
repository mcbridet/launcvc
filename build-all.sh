#!/bin/bash
set -e

GOOS=linux GOARCH=amd64 go build -o bin/launcvc_linux_amd64 main.go
GOOS=linux GOARCH=arm go build -o bin/launcvc_linux_arm main.go
GOOS=linux GOARCH=arm64 go build -o bin/launcvc_linux_arm64 main.go

GOOS=windows GOARCH=amd64 go build -o bin/launcvc_win_amd64.exe main.go

GOOS=darwin GOARCH=amd64 go build -o bin/launcvc_darwin_amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o bin/launcvc_darwin_arm64 main.go
