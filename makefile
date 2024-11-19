# Makefile for Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=plan.exe
MAIN_PATH=./cmd/main.go

# Build the project
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(BINARY_NAME)

all: build run

build:
	mkdir $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)


clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

run:
	./$(BINARY_PATH)


