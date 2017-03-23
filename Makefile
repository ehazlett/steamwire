CGO_ENABLED=0
GOOS?=linux
GOARCH?=amd64
COMMIT=`git rev-parse --short HEAD`
APP=steamwire
OUTPUT?=$(APP)
REPO?=ehazlett/$(APP)
TAG?=latest
DEPS=$(shell go list ./... | grep -v /vendor/)

all: build

build:
	@cd cmd/$(APP) && go build -o $(OUTPUT) -ldflags "-w -X github.com/ehazlett/$(APP)/version.GitCommit=$(COMMIT)" .

build-static:
	@cd cmd/$(APP) && go build -o $(OUTPUT) -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/ehazlett/$(APP)/version.GitCommit=$(COMMIT)" .

test: check
	@go test -v $(DEPS)

check:
	@go vet -v $(DEPS)
	@golint $(DEPS)

release:
	@script/release.sh

clean:
	@rm -rf cmd/$(APP)/$(APP)
	@rm -rf build

.PHONY: add-deps build build-static test clean
