GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

all: deps build run
build:
	$(GOBUILD) -o assignment1
deps:
	$(GOGET) golang.org/x/sys/unix
run:
	sudo ./assignment1

