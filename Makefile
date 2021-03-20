MODULE := github.com/spotify-bot/server

deps:
	go mod download

build: deps
	env GOOS=linux GOARCH=amd64 go build -o build/webserver cmd/main.go

docker:
	docker build -t webserver .

lint:
	golangci-lint run --disable errcheck
