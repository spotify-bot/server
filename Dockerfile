FROM golang:1.15 as builder

ARG CMD
ARG ARCH=amd64
ARG OS=linux

SHELL ["/bin/bash", "-c"]

ENV CGO_ENABLED 0
ENV HOME /app

WORKDIR /app
COPY . .
RUN go mod download
RUN	env GOOS=$OS GOARCH=$ARCH go build -o build/$CMD github.com/spotify-bot/server/cmd/$CMD

FROM alpine:3.12 as app

WORKDIR /app
COPY --from=builder /app/build/$CMD /app
