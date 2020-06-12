GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=chat
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
server: 
	go run server.go
client: 
	go run client.go
build: 
	go build -v -o $(BINARY_NAME) server.go
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v