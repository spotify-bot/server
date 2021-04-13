#!make

SHELL=/bin/sh

SERVICE_COMPOSE_FILE ?= docker-compose-service.yaml
CONFIF_ENV_FILE ?= app.env

deps:
	go mod download

build: deps
	env GOOS=linux GOARCH=amd64 go build -o build cmd/main.go

docker:
	docker build -t webserver .

lint:
	golangci-lint run --disable errcheck
service.build:
	docker compose -f ${SERVICE_COMPOSE_FILE} build --pull

service.run:
	docker compose -f ${SERVICE_COMPOSE_FILE} \
	--env-file ${CONFIF_ENV_FILE} \
	up -d
