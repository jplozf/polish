#!/bin/sh

# Linux Binaries
GOOS=linux GOARCH=amd64 go build -o bin/linux/polish-amd64
GOOS=linux GOARCH=386 go build -o bin/linux/polish-386

# Windows Binaries
GOOS=windows GOARCH=amd64 go build -o bin/windows/polish-amd64.exe
GOOS=windows GOARCH=386 go build -o bin/windows/polish-386.exe

# macOS Binaries
GOOS=darwin GOARCH=amd64 go build -o bin/macOS/polish-amd64
GOOS=darwin GOARCH=arm64 go build -o bin/macOS/polish-arm64
