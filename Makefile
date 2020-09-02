 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

all: build
build:
	$(GOGET) github.com/PuerkitoBio/goquery
	$(GOBUILD) src/*
clean:
	$(GOCLEAN)
	rm -f phrase_indexer
