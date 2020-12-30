package main

import (
	"context"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/koskalak/mamal/internal/mongo"
	"github.com/koskalak/mamal/internal/webserver"
)

type config struct {
	Env                 string `env:"ENV" envDefault:"development"`
	Addr                string `env:"ADDR" envDefault:":8000"`
	JWTSecret           string `env:"JWT_SECRET,required"`
	MongoDSN            string `env:"MONGODB_DSN" envDefault:"mongodb://localhost:27017/mamal"`
	SpotifyClientID     string `env:"SPOTIFY_CLIENT_ID"`
	SpotifyClientSecret string `env:"SPOTIFY_CLIENT_SECRET"`
}

func main() {
	cfg := newConfig()
	ctx := context.Background()

	mongoStorage, err := mongo.NewMongoStorage(ctx, mongo.MongoStorageOptions{
		DSN: cfg.MongoDSN,
	})

	w := webserver.New(webserver.WebServerOptions{
		Mongo:               mongoStorage,
		SpotifyClientID:     cfg.SpotifyClientID,
		SpotifyClientSecret: cfg.SpotifyClientSecret,
	})

	err = w.Start(cfg.Addr)

	if err != nil {
		log.Fatalln("Cannot start server", err)
	}
}

func newConfig() *config {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Unable to load .env %v\n", err)
	}

	cfg := &config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Unable to parse env vars %v\n", err)
	}

	return cfg
}
