package telegram

import (
	"log"

	"github.com/koskalak/mamal/config"
	tgbotapi "github.com/mohammadkarimi23/telegram-bot-api/v5"
	"strconv"
)

func Start(token string) {
	bot, err := tgbotapi.NewBotAPI(token) //FIXME change to use configs
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	processMessage(bot)
}

func processMessage(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s]\n %s\n*********", update.Message.From.UserName, update.Message.Text)

		txt := "Please use the following link for auth: \n" + config.AppConfig.Webserver.Address + "/auth/telegram?user_id=" + strconv.Itoa(update.Message.From.ID) //FIXME change to user_ID
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func getMessage(user_id string)                   {} //Get link for song from spotify.
func shareMessage(message string, chat_id string) {} //share message to recipient.
