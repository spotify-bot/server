package telegram

import (
	"log"

	"github.com/koskalak/mamal/config"
	tgbotapi "github.com/mohammadkarimi23/telegram-bot-api/v5"
	"strconv"
)

type TGBot struct {
	bot *tgbotapi.BotAPI
}

func New(token string) *TGBot {
	bot, err := tgbotapi.NewBotAPI(token) //FIXME change to use configs
	if err != nil {
		log.Panic(err)
	}
	telegramBot := TGBot{
		bot: bot,
	}
	return &telegramBot
}

func (tb *TGBot) Start() {
	tb.bot.Debug = true

	log.Printf("Authorized on account %s", tb.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := tb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			tb.processDirectMessage(&update)
		} else if update.InlineQuery != nil {
			//tb.processInlineQuery(&update)
		} else {
			continue
		}
	}
}

func (tb *TGBot) processDirectMessage(update *tgbotapi.Update) {
	log.Printf("[%s]\n %s\n*********", update.Message.From.UserName, update.Message.Text)

	txt := "Please use the following link for auth: \n" + config.AppConfig.Webserver.Address + "/auth/telegram?user_id=" + strconv.Itoa(update.Message.From.ID) //FIXME change to user_ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
	msg.ReplyToMessageID = update.Message.MessageID

	tb.bot.Send(msg)
}

func getMessage(user_id string)                   {} //Get link for song from spotify.
func shareMessage(message string, chat_id string) {} //share message to recipient.
