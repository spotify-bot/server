#!make

SHELL=/bin/sh

CONFIF_ENV_FILE ?= app.env

deps:
	go mod download

build: deps
	env GOOS=linux GOARCH=amd64 go build -o build cmd/main.go

docker:
	docker build -t webserver .

lint:
	golangci-lint run --disable errcheck

run:
	docker compose \
    --env-file ${CONFIF_ENV_FILE} \
    up -d
