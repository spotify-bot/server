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
	log.Printf("Authorized on account %s", bot.Self.UserName)
	telegramBot := TGBot{
		bot: bot,
	}
	return &telegramBot
}

func (tb *TGBot) Start() {
	tb.bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := tb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.IsCommand() {
			tb.processCommand(&update)
		} else if update.Message != nil {
			tb.processDirectMessage(&update)
		} else if update.InlineQuery != nil {
			tb.processInlineQuery(&update)
		} else {
			continue
		}
	}
}

func (tb *TGBot) processCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	//	txt := "Please use the following link for auth: \n" + //FIXME change to user_ID
	switch update.Message.Command() {
	case "help":
		msg.Text = "type /auth to authenticate to Spotify."
	case "auth":
		msg.ReplyMarkup = getAuthMessage(strconv.Itoa(update.Message.From.ID))
		msg.Text = "Click here to login to Spotify."
	default:
		msg.Text = "I don't know that command"
	}
	tb.bot.Send(msg)
}

func (tb *TGBot) processDirectMessage(update *tgbotapi.Update) {
	log.Printf("Message from [%s]:  %s\n", update.Message.From.UserName, update.Message.Text)

	txt := "Please use the following link for auth: \n" + config.AppConfig.Webserver.Address + "/auth/telegram?user_id=" + strconv.Itoa(update.Message.From.ID) //FIXME change to user_ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)

	tb.bot.Send(msg)
}

func (tb *TGBot) processInlineQuery(update *tgbotapi.Update) {

	log.Printf("Inline Queryyyyyyyyyyyyyyyy")

	article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Echo", update.InlineQuery.Query)
	article.Description = update.InlineQuery.Query

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{article},
	}

	if _, err := tb.bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println(err)
	}
}

func getAuthMessage(userID string) tgbotapi.InlineKeyboardMarkup {

	link := "http://" + config.AppConfig.Webserver.Address + "/auth/telegram?user_id=" + userID //FIXME dev config
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Login to Spotify", link),
		),
	)
}

func getMessage(user_id string)                   {} //Get link for song from spotify.
func shareMessage(message string, chat_id string) {} //share message to recipient.
