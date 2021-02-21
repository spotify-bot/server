MODULE := github.com/spotify-bot/server
TELEGRAM_CMD := $(MODULE)/cmd/telegram
WEBSERVER_CMD := $(MODULE)/cmd/webserver

deps:
	go mod download

build: build.telegram build.webserver

build.telegram: deps
	env GOOS=linux GOARCH=amd64 go build -o build/telegram $(TELEGRAM_CMD)

build.webserver: deps
	env GOOS=linux GOARCH=amd64 go build -o build/webserver $(WEBSERVER_CMD)

docker.telegram:
	docker build -t telegram --build-arg CMD=telegram . 

docker.webserver:
	docker build -t webserver --build-arg CMD=webserver .
