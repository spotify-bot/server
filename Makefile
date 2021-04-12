#!make

SHELL=/bin/sh
SERVICE_COMPOSE_FILE=docker-compose-service.yaml

-include: .env .env.local .env.*.local

deps:
	go mod download

build: deps
	env GOOS=linux GOARCH=amd64 go build -o build cmd/main.go

docker:
	docker build -t webserver .

lint:
	golangci-lint run --disable errcheck

service.build:
	docker-compose -f ${SERVICE_COMPOSE_FILE} build --pull

service.run:
	ADDRESS=${ADDRESS} \
	MONGO_DSN=${MONGO_DSN} \
	API_SERVER_ADDRESS=${API_SERVER_ADDRESS} \
	CLIENT_ID=${CLIENT_ID} \
	CLIENT_SECRET=${CLIENT_SECRET} \
	docker-compose -f ${SERVICE_COMPOSE_FILE} up -d
