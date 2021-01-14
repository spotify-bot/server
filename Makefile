MODULE := github.com/koskalak/mamal
TELEGRAM_CMD := $(MODULE)/cmd/telegram
WEBSERVER_CMD := $(MODULE)/cmd/webserver

deps:
	go mod download

build: build.telegram build.webserver

build.telegram:
	env GOOS=linux GOARCH=amd64 go build -o build/telegram $(TELEGRAM_CMD)

build.webserver:
	env GOOS=linux GOARCH=amd64 go build -o build/webserver $(WEBSERVER_CMD)

build.docker: build
	docker build -t telegram --build-arg BINARY=telegram . 
	docker build -t webserver --build-arg BINARY=webserver .


