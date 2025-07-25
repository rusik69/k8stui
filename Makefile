# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=k8stui
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the project
all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/$(BINARY_NAME)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/$(BINARY_NAME)
	./$(BINARY_NAME)

deps:
	$(GOMOD) tidy
	$(GOMOD) download

# Cross compilation
build-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-darwin-arm64 -v ./cmd/$(BINARY_NAME)

build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 -v ./cmd/$(BINARY_NAME)

build-all: build build-darwin-arm64 build-linux-amd64

# Docker
.PHONY: all build test clean run deps build-linux
