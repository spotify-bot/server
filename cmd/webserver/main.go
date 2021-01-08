package main

import (
	"context"
	"log"

	"github.com/koskalak/mamal/internal/config"
	"github.com/koskalak/mamal/internal/spotify"
	"github.com/koskalak/mamal/internal/webserver"
	"golang.org/x/oauth2"
	spotifyOauth "golang.org/x/oauth2/spotify"
)

func main() {
	ctx := context.Background()

	authConf := &oauth2.Config{
		ClientID:     config.AppConfig.Spotify.SpotifyClientID,
		ClientSecret: config.AppConfig.Spotify.SpotifyClientSecret,
		Scopes:       []string{"user-read-playback-state"},
		Endpoint:     spotifyOauth.Endpoint,
		RedirectURL:  "http://" + config.AppConfig.Webserver.Address + "/auth/callback", //FIXME
	}

	s, err := spotify.New(ctx, spotify.ProviderOptions{
		DatabaseDSN: config.AppConfig.Webserver.MongoDSN,
		AuthConfig:  authConf,
	})
	if err != nil {
		log.Fatal("Failed to connect to Mongo Storage", err)
	}

	w := webserver.New(webserver.WebServerOptions{
		Spotify: s,
	})

	err = w.Start(config.AppConfig.Webserver.Address)

	if err != nil {
		log.Fatalln("Cannot start server", err)
	}
}
