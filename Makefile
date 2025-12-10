APP_NAME=workout_app
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: build run test clean fmt lint

build:
	go build -o $(APP_NAME) ./...

dev:
	air

test:
	go test ./...

fmt:
	gofmt -s -w $(GO_FILES)

lint:
	golint ./...

clean:
	rm -f $(APP_NAME)

swag:
	swag init
