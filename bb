#!/bin/sh

export GIT_COMMIT=$(git rev-parse --short HEAD)
export BUILD_TIME=`date '+%Y%m%d-%H%M%S'`

# Linux Binaries
GOOS=linux GOARCH=amd64 go build -v -ldflags="-X 'main.GIT_COMMIT=$GIT_COMMIT' -X 'main.BUILD_TIME=$BUILD_TIME'" -o bin/linux/polish-amd64
GOOS=linux GOARCH=386 go build -v -ldflags="-X 'main.GIT_COMMIT=$GIT_COMMIT' -X 'main.BUILD_TIME=$BUILD_TIME'" -o bin/linux/polish-386

# Windows Binaries
GOOS=windows GOARCH=amd64 go build -v -ldflags="-X 'main.GIT_COMMIT=$GIT_COMMIT' -X 'main.BUILD_TIME=$BUILD_TIME'" -o bin/windows/polish-amd64.exe
GOOS=windows GOARCH=386 go build -v -ldflags="-X 'main.GIT_COMMIT=$GIT_COMMIT' -X 'main.BUILD_TIME=$BUILD_TIME'" -o bin/windows/polish-386.exe

# macOS Binaries
GOOS=darwin GOARCH=amd64 go build -v -ldflags="-X 'main.GIT_COMMIT=$GIT_COMMIT' -X 'main.BUILD_TIME=$BUILD_TIME'" -o bin/macOS/polish-amd64
GOOS=darwin GOARCH=arm64 go build -v -ldflags="-X 'main.GIT_COMMIT=$GIT_COMMIT' -X 'main.BUILD_TIME=$BUILD_TIME'" -o bin/macOS/polish-arm64
