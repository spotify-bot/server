package main

import (
	"context"
	"log"

	"github.com/koskalak/mamal/internal/config"
	"github.com/koskalak/mamal/internal/mongo"
	"github.com/koskalak/mamal/internal/webserver"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

func main() {
	ctx := context.Background()

	mongoStorage, err := mongo.NewMongoStorage(ctx, mongo.MongoStorageOptions{
		DSN: config.AppConfig.Webserver.MongoDSN,
	})

	authConf := &oauth2.Config{
		ClientID:     config.AppConfig.Spotify.SpotifyClientID,
		ClientSecret: config.AppConfig.Spotify.SpotifyClientSecret,
		Scopes:       []string{"user-read-playback-state"},
		Endpoint:     spotify.Endpoint,
		RedirectURL:  "http://" + config.AppConfig.Webserver.Address + "/auth/callback", //FIXME
	}

	w := webserver.New(webserver.WebServerOptions{
		Mongo:      mongoStorage,
		AuthConfig: authConf,
	})

	err = w.Start(config.AppConfig.Webserver.Address)

	if err != nil {
		log.Fatalln("Cannot start server", err)
	}
}
