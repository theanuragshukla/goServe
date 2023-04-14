#!/usr/bin/bash

set -e

go get ./...
# Ensure that the build directory exists
mkdir -p build
# Build for Mac
echo "Building for Mac"


GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/gin-gonic/gin.Mode=release'" -o build/goserve-mac ./src/main.go

# Build for Linux
echo "Building for Linux"
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/gin-gonic/gin.Mode=release'" -o build/goserve-linux ./src/main.go

# Build for Windows
echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/gin-gonic/gin.Mode=release'" -o build/goserve-windows.exe ./src/main.go

