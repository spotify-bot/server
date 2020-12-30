package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/koskalak/mamal/internal/webserver"
)

type config struct {
	Env       string `env:"ENV" envDefault:"development"`
	Addr      string `env:"ADDR" envDefault:":8000"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

func main() {
	cfg := newConfig()

	w := webserver.New()

	err := w.Start(cfg.Addr)

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
