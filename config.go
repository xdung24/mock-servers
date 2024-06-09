package main

import "github.com/spf13/viper"

type Config struct {
	DataFolder  string `mapstructure:"DATA_FOLDER"`
	UseFsNotify bool   `mapstructure:"USE_FSNOTIFY"`
	UsePolling  bool   `mapstructure:"USE_POLLING"`
	PollingTime int    `mapstructure:"POLLING_TIME"`
}

func getConfig() Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	dataFolder := viper.GetString("DATA_FOLDER")
	if dataFolder == "" {
		dataFolder = "./data"
	}

	useFsNotify := viper.GetBool("USE_FSNOTIFY")
	usePolling := viper.GetBool("USE_POLLING")
	pollingTime := viper.GetInt("POLLING_TIME")

	return Config{
		DataFolder:  dataFolder,
		UseFsNotify: useFsNotify,
		UsePolling:  usePolling,
		PollingTime: pollingTime,
	}
}
