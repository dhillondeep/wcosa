# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=wio
BINARY_UNIX=$(BINARY_NAME)_unix

all: build run

get:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/davecgh/go-spew/spew
	go get -u github.com/stretchr/testify
	go get -u github.com/pmezard/go-difflib/difflib
	go get -u github.com/kevinburke/go-bindata/...

build:
	@echo Building $(BINARY_NAME) project:
	@cd "$(CURDIR)/pkg/util/sys" && go-bindata -nomemcopy -pkg sys -prefix ../../../ ../../../assets/...
	@cd "$(CURDIR)/cmd/$(BINARY_NAME)" && $(GOBUILD) -ldflags="-s -w" -o ../../bin/$(BINARY_NAME) -v
	@echo Done!

clean:
	@echo Cleaning $(BINARY_NAME) project files:
	@$(GOCLEAN)
	@rm -f bin/$(BINARY_NAME)
	@rm -f bin/$(BINARY_UNIX)
	@echo Done!

run:
	@./bin/$(BINARY_NAME) ${ARGS}
