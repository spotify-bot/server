package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/koskalak/mamal/internal/platform/telegram"
)

type config struct {
	telegramBotID string `env:"TELEGRAM_BOT_ID,required"`
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

func main() {
	//cfg := newConfig()
	telegram.Start(os.Getenv("API_KEY")) //FIXME change to use configs
}
