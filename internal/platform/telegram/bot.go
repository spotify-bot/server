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

	txt := "Please use the following link for auth: \n" + config.AppConfig.Spotify.ApiServerAddress + "/auth/telegram?user_id=" + strconv.Itoa(update.Message.From.ID) //FIXME change to user_ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)

	tb.bot.Send(msg)
}

func (tb *TGBot) processInlineQuery(update *tgbotapi.Update) {

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
	}

	rp, err := tb.spotify.GetRecentlyPlayed(spotify.PlatformTelegram, strconv.Itoa(update.InlineQuery.From.ID))
	if err != nil { //If not logged in, show the log in keyboard button
		log.Println("Failed to get recently played song", err)
		inlineConf.SwitchPMText = "Login to Spotify"
		inlineConf.SwitchPMParameter = "auth"
	} else {
		songLink := spotify.OpenSpotifyTrackEndpoint + rp.ID
		article := getTrackQueryResult(update.InlineQuery.ID, rp.Name, songLink, songLink)
		inlineConf.Results = []interface{}{article}
	}

	if _, err := tb.bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println("Failed to answer inline query: ", err)
	}
}
