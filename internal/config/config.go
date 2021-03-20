package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

var (
	ConfigPath string
)

func init() {
	// Register flags
	flag.StringVar(&ConfigPath, "config-path", "", "Path to config directory")
	flag.Parse()

	if ConfigPath == "" {
		panic("No config path provided to run the application")
	}

	initConfig()
}

func loadConfig(configStruct interface{}) {
	viperInstance := viper.New()
	viperInstance.SetConfigName(ConfigPath)
	viperInstance.AddConfigPath(".")
	viperInstance.SetConfigType("yaml")

	err := viperInstance.ReadInConfig()
	if err != nil {
		log.Panic("failed to read Config, ", err)
	}

	err = viperInstance.Unmarshal(configStruct)
	if err != nil {
		log.Panic("failed to unmarshal Config file")
	}
	log.Printf("loaded config: %v\n", ConfigPath)
}
