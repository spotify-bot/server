package config

type Config struct {
	Address          string `mapstructure:"ADDRESS"`
	MongoDSN         string `mapstructure:"MONGO_DSN"`
	APIServerAddress string `mapstructure:"API_SERVER_ADDRESS"`
	ClientID         string `mapstructure:"CLIENT_ID"`
	ClientSecret     string `mapstructure:"CLIENT_SECRET"`
}

var AppConfig Config

func initConfig() {
	loadConfig(&AppConfig)
}
