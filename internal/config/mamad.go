package config

type Config struct {
	Webserver   Webserver   `mapstructure:"webserver"`
	TelegramBot TelegramBot `mapstructure:"telegramBot"`
	Spotify     Spotify     `mapstructure:"spotify"`
}

type Webserver struct {
	Address   string `mapstructure:"address"`
	JWTSecret string `mapstructure:"jwtSecret"`
	MongoDSN  string `mapstructure:"mongoDSN"`
}

type TelegramBot struct {
	APIToken string `mapstructure:"apiToken"`
}

type Spotify struct {
	ApiServerAddress    string `mapstructure:"apiServerAddress"`
	SpotifyClientID     string `mapstructure:"clientID"`
	SpotifyClientSecret string `mapstructure:"clientSecret"`
}

var AppConfig Config

func initMamalConfig() {
	loadConfig(&AppConfig)
}
