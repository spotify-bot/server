package telegram

import (
	"github.com/koskalak/mamal/internal/config"
	tgbotapi "github.com/mohammadkarimi23/telegram-bot-api/v5"
)

func getAuthMessage(userID string) tgbotapi.InlineKeyboardMarkup {

	link := "http://" + config.AppConfig.Webserver.Address + "/auth/telegram?user_id=" + userID //FIXME dev config
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Login to Spotify", link),
		),
	)
}
