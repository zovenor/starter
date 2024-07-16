VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)

build:
	@go build -ldflags "-X main.version=$(VERSION)" -o ./bin/starter ./app/starter.go

run:
	@go run  -ldflags "-X main.version=$(VERSION)" ./app/starter.go

install:
	@go install -ldflags "-X main.version=$(VERSION)" ./app/starter.go 

.PHONY: run install
