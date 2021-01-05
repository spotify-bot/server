package config

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	initMamalConfig()
}

func loadConfig(configStruct interface{}) {
	viperInstance := viper.New()
	viperInstance.SetConfigName(ConfigPath)
	viperInstance.AddConfigPath(".")
	viperInstance.SetConfigType("yaml")

	err := viperInstance.ReadInConfig()
	if err != nil {
		errorString := fmt.Sprintln("failed to read '%v' (err: %v)\n", ConfigPath, err)
		panic(errorString)
	}

	err = viperInstance.Unmarshal(configStruct)
	if err != nil {
		errorString := fmt.Sprintln("failed to unmarshal '%v' (err: %v)\n", ConfigPath, err)
		panic(errorString)
	}
	fmt.Printf("loaded config: %v\n", ConfigPath)
}
