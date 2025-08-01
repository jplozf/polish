#!/bin/bash

# Get the major version from main.go
MAJOR_VERSION=$(grep 'const majorVersion' main.go | awk -F'"' '{print $2}')

# Get the short commit hash
GIT_HASH=$(git rev-parse --short HEAD)

# Get the total commit count
COMMIT_COUNT=$(git rev-list --count HEAD)

# Combine them to form the version string (major.commits-hash)
VERSION="${MAJOR_VERSION}.${COMMIT_COUNT}-${GIT_HASH}"

echo "Building Polish interpreter with version: ${VERSION}"

# Build the Go application, embedding the version string
go build -a -ldflags "-X main.version=${VERSION}" -o polish .

if [ $? -eq 0 ]; then
    echo "Build successful: ./polish"
else
    echo "Build failed."
    exit 1
fi