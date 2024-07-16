build:
	@go build -o ./bin/starter ./app/starter.go

run:
	@go run ./app/starter.go

install:
	@go install ./app/starter.go

.PHONY: run install
