package telegram

import (
	"log"

	"github.com/koskalak/mamal/internal/config"
	"github.com/koskalak/mamal/internal/spotify"
	tgbotapi "github.com/mohammadkarimi23/telegram-bot-api/v5"
	"strconv"
)

type TGBot struct {
	bot     *tgbotapi.BotAPI
	spotify *spotify.SpotifyProvider
}

type TGBotOptions struct {
	Token           string
	SpotifyProvider *spotify.SpotifyProvider
}

func New(opts TGBotOptions) *TGBot {
	bot, err := tgbotapi.NewBotAPI(opts.Token) //FIXME change to use configs
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	telegramBot := TGBot{
		bot:     bot,
		spotify: opts.SpotifyProvider,
	}
	return &telegramBot
}

func (tb *TGBot) Start() {
	tb.bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := tb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery != nil {
			tb.processInlineQuery(&update)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				tb.processCommand(&update)
			} else {
				tb.processDirectMessage(&update)
			}
		} else {
			continue
		}
	}
}

func (tb *TGBot) processCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	//	txt := "Please use the following link for auth: \n" + //FIXME change to user_ID
	switch update.Message.Command() {
	case "start":
		msg.ReplyMarkup = getAuthMessage(strconv.Itoa(update.Message.From.ID))
		msg.Text = "Click here to login to Spotify."
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

	rp, err := tb.spotify.GetRecentlyPlayed("telegram", strconv.Itoa(update.InlineQuery.From.ID))
	if err != nil {
		log.Println("Failed to get recently played song", err)
	} else {
		log.Println("Here comes the result: ", spotify.OpenSpotifyTrackEndpoint+rp.ID)
	}

	article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Echo", update.InlineQuery.Query)
	article.Description = update.InlineQuery.Query

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID:     update.InlineQuery.ID,
		IsPersonal:        true,
		CacheTime:         0,
		Results:           []interface{}{article},
		SwitchPMText:      "Login to Spotify",
		SwitchPMParameter: "auth",
	}

	if _, err := tb.bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println(err)
	}
}

func getMessage(user_id string)                   {} //Get link for song from spotify.
func shareMessage(message string, chat_id string) {} //share message to recipient.
