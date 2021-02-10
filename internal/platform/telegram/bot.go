package telegram

import (
	"log"

	"github.com/koskalak/mamal/internal/config"
	"github.com/koskalak/mamal/internal/spotify"
	tgbotapi "github.com/mohammadkarimi23/telegram-bot-api/v5"
	"strconv"
	"strings"
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
		} else if update.CallbackQuery != nil {
			tb.processCallbackQuery(&update)
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
	if _, err := tb.bot.Send(msg); err != nil {
		log.Fatal("Failed to send message ", err)
	}
}

func (tb *TGBot) processDirectMessage(update *tgbotapi.Update) {
	log.Printf("Message from [%s]:  %s\n", update.Message.From.UserName, update.Message.Text)

	txt := "Please use the following link for auth: \n" + config.AppConfig.Spotify.ApiServerAddress + "/auth/telegram?user_id=" + strconv.Itoa(update.Message.From.ID) //FIXME change to user_ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)

	if _, err := tb.bot.Send(msg); err != nil {
		log.Fatal("Failed to send message ", err)
	}
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
		article := getTrackQueryResult(update.InlineQuery.ID, rp)
		inlineConf.Results = []interface{}{article}
	}

	if _, err := tb.bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println("Failed to answer inline query: ", err)
	}
}

func (tb *TGBot) processCallbackQuery(update *tgbotapi.Update) {

	if !strings.Contains(update.CallbackQuery.Data, "#") {
		log.Println("Unknown callback format")
		return //TODO should send proper response to sender
	}
	splittedQuery := strings.SplitAfterN(update.CallbackQuery.Data, "#", 2)
	trackURI := spotify.TrackURIPrefix + splittedQuery[1]
	queryType, err := strconv.Atoi(strings.Trim(splittedQuery[0], "#"))
	if err != nil {
		log.Println("Unknown callback format:")
		return //TODO should send proper response to sender
	}
	callbackMessage := ""
	callbackConf := tgbotapi.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
	}

	switch queryType {
	case 1:
		err := tb.spotify.PlaySong(spotify.PlatformTelegram, strconv.Itoa(update.CallbackQuery.From.ID), trackURI)
		if err != nil {
			callbackMessage = "Failed to play song"
			log.Println("Failed to play song: ", err)
		} else {
			callbackMessage = "Enjoy!"
		}
	case 2:
		err := tb.spotify.AddSongToQueue(spotify.PlatformTelegram, strconv.Itoa(update.CallbackQuery.From.ID), trackURI)
		if err != nil {
			callbackMessage = "Failed to add song to queue"
			log.Println("Failed to add song to queue: ", err)
		} else {
			callbackMessage = "Song Added To Queue"
		}
	default:
		log.Println("Callback Query type unknown")
	}
	callbackConf.Text = callbackMessage

	if _, err := tb.bot.AnswerCallbackQuery(callbackConf); err != nil {
		log.Println("Failed to answer inline query: ", err)
	}
}
