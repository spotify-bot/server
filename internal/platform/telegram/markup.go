package telegram

import (
	"github.com/koskalak/mamal/internal/config"
	"github.com/koskalak/mamal/internal/spotify"
	tgbotapi "github.com/mohammadkarimi23/telegram-bot-api/v5"
)

func getAuthMessage(userID string) tgbotapi.InlineKeyboardMarkup {

	link := "http://" + config.AppConfig.Spotify.ApiServerAddress + "/auth/telegram?user_id=" + userID //FIXME dev config
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Login to Spotify", link),
		),
	)
}

func getTrackQueryResult(id string, track *spotify.Item) tgbotapi.InlineQueryResultArticle {
	trackLink := spotify.OpenSpotifyTrackEndpoint + track.ID
	return tgbotapi.InlineQueryResultArticle{
		Type:  "article",
		ID:    id,
		Title: track.Name,
		URL:   trackLink,
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text: trackLink,
		},
		ThumbURL:    track.Album.Images[1].URL, //FIXME use smallest image
		ThumbWidth:  track.Album.Images[1].Width,
		ThumbHeight: track.Album.Images[1].Height,
		Description: track.Album.Artists[0].Name, //FIXME
	}
}
