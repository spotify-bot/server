package config

type Config struct {
	Webserver   Webserver   `mapstructure:"webserver"`
	TelegramBot TelegramBot `mapstructure:"telegramBot"`
	Spotify     Spotify     `mapstructure:"spotify"`
}

type Webserver struct {
	Address   string `mapstructure:"address"`
	JWTSecret string `mapstructure:"jwtSecret"`
	MongoSDN  string `mapstructure:"mongoSDN"`
}

type TelegramBot struct {
	APIToken string `mapstructure:"apiToken"`
}

type Spotify struct {
	SpotifyClientID     string `mapstructure:"clientID"`
	SpotifyClientSecret string `mapstructure:"clientSecret"`
}

var AppConfig Config

func initMamalConfig() {
	loadConfig(&AppConfig)
}
