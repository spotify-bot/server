package main

import (
	"github.com/koskalak/mamal/internal/config"
	"github.com/koskalak/mamal/internal/platform/telegram"
)

func main() {
	tgbot := telegram.New(config.AppConfig.TelegramBot.APIToken)
	tgbot.Start()
}
