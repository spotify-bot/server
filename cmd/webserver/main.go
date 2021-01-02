package main

import (
	"context"
	"log"

	"github.com/koskalak/mamal/config"
	"github.com/koskalak/mamal/internal/mongo"
	"github.com/koskalak/mamal/internal/webserver"
)

func main() {
	ctx := context.Background()

	mongoStorage, err := mongo.NewMongoStorage(ctx, mongo.MongoStorageOptions{
		DSN: config.AppConfig.Webserver.MongoSDN,
	})

	w := webserver.New(webserver.WebServerOptions{
		Mongo:               mongoStorage,
		SpotifyClientID:     config.AppConfig.Spotify.SpotifyClientID,
		SpotifyClientSecret: config.AppConfig.Spotify.SpotifyClientSecret,
	})

	err = w.Start(config.AppConfig.Webserver.Address)

	if err != nil {
		log.Fatalln("Cannot start server", err)
	}
}
