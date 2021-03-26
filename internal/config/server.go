package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Address          string `mapstructure:"ADDRESS"`
	MongoDSN         string `mapstructure:"MONGO_DSN"`
	APIServerAddress string `mapstructure:"API_SERVER_ADDRESS"`
	ClientID         string `mapstructure:"CLIENT_ID"`
	ClientSecret     string `mapstructure:"CLIENT_SECRET"`
}

var (
	ConfigName string
	AppConfig  Config
)

func init() {
	// Register flags
	flag.StringVar(&ConfigName, "config-name", "", "Path to config directory")
	flag.Parse()

	loadConfig(&AppConfig)
}

func loadConfig(configStruct interface{}) {
	viperInstance := viper.New()

	// If Configname is passed, read the config from file
	if ConfigName != "" {
		viperInstance.SetConfigName(ConfigName)
		viperInstance.AddConfigPath(".")
		viperInstance.SetConfigType("env")

		err := viperInstance.ReadInConfig()
		if err != nil {
			log.Panic("Config file is not set", err)
		}

		err = viperInstance.Unmarshal(configStruct)
		if err != nil {
			log.Panic("Failed to Unmarshal Config file: ", err)
		}
		log.Printf("loaded config: %v\n", ConfigName)
	} else {
		viperInstance.AutomaticEnv()
		address := viperInstance.GetString("ADDRESS")
		mongoDSN := viperInstance.GetString("MONGO_DSN")
		apiServerAddress := viperInstance.GetString("API_SERVER_ADDRESS")
		clientID := viperInstance.GetString("CLIENT_ID")
		clientSecret := viperInstance.GetString("CLIENT_SECRET")
		if len(address) > 0 && len(mongoDSN) > 0 && len(apiServerAddress) > 0 && len(clientID) > 0 && len(clientSecret) > 0 {
			AppConfig = Config{
				Address:          address,
				MongoDSN:         mongoDSN,
				APIServerAddress: apiServerAddress,
				ClientID:         clientID,
				ClientSecret:     clientSecret,
			}
		} else {
			log.Panic("ENV Vars should be set")
		}
	}
}
