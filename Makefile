 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=phrase_indexer

all: build
build:
	$(GOGET) github.com/PuerkitoBio/goquery
	$(GOBUILD) -o $(BINARY_NAME) src/*
clean:
	$(GOCLEAN)
	rm -f phrase_indexer
