package main

import (
	"github.com/koskalak/mamal/config"
	"github.com/koskalak/mamal/internal/platform/telegram"
)

func main() {
	telegram.Start(config.AppConfig.TelegramBot.APIToken)
}
