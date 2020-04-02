 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

all: build
build:
	$(GOBUILD) src/*
clean:
	$(GOCLEAN)
	rm -f phrase_indexer
