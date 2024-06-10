package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	DockerRunning bool   `mapstructure:"DOCKER_RUNNING"`
	DataFolder    string `mapstructure:"DATA_FOLDER"`
	UseFsNotify   bool   `mapstructure:"USE_FSNOTIFY"`
	UsePolling    bool   `mapstructure:"USE_POLLING"`
	PollingTime   int    `mapstructure:"POLLING_TIME"`
	WebEngine     string `mapstructure:"WEB_ENGINE"`
}

// Get the configuration from order:
// - .env file
// - the environment variables
// - the command line flags
func getEnvConfig() Config {
	// Initialize Viper
	viper.SetConfigFile(".env") // Load .env file
	viper.AutomaticEnv()        // Read environment variables

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading the config file: ", err)
	}

	// Define a Cobra command
	rootCmd := &cobra.Command{
		Use:   "MockServer",
		Short: "MockServer is an application to serve http servers with predefined responses.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// Add flags to the command
	rootCmd.Flags().String("data-folder", "./data/", "Location to data folder")
	rootCmd.Flags().Bool("use-fsnotify", false, "Use FsNotify to watch for changes in the data folder")
	rootCmd.Flags().Bool("use-polling", false, "Use Polling to watch for changes in the data folder (only use this if fsnotify is not working)")
	rootCmd.Flags().Int("polling-time", 10, "Polling time in seconds")
	rootCmd.Flags().String("web-engine", "", "Web engine to use (gorilla, gin, echo, fiber)")

	// Bind flags to Viper keys
	viper.BindPFlag("DATA_FOLDER", rootCmd.Flags().Lookup("data-folder"))
	viper.BindPFlag("USE_FSNOTIFY", rootCmd.Flags().Lookup("use-fsnotify"))
	viper.BindPFlag("USE_POLLING", rootCmd.Flags().Lookup("use-polling"))
	viper.BindPFlag("POLLING_TIME", rootCmd.Flags().Lookup("polling-time"))
	viper.BindPFlag("WEB_ENGINE", rootCmd.Flags().Lookup("web-engine"))

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return Config{
		DockerRunning: viper.GetBool("DOCKER_RUNNING"),
		DataFolder:    viper.GetString("DATA_FOLDER"),
		UseFsNotify:   viper.GetBool("USE_FSNOTIFY"),
		UsePolling:    viper.GetBool("USE_POLLING"),
		PollingTime:   viper.GetInt("POLLING_TIME"),
		WebEngine:     viper.GetString("WEB_ENGINE"),
	}
}

func (config Config) Validate() error {
	if config.DataFolder == "" {
		return errors.New("data folder is required")
	}

	// Check if data folder exists
	if _, err := os.Stat(config.DataFolder); os.IsNotExist(err) {
		return errors.New("data folder does not exist")
	}

	// Can not use both fsnotify and polling
	if config.UseFsNotify && config.UsePolling {
		return errors.New("can not use both fsnotify and polling")
	}

	// Polling time should be greater than 0s
	if config.UsePolling && config.PollingTime <= 0 {
		return errors.New("polling time should be greater than 0")
	}

	// Web engine should be gin, gorilla, echo, fiber
	engines := []string{"gin", "gorilla", "echo", "fiber"}
	if !slices.Contains(engines, config.WebEngine) {
		return errors.New("web engine should be gin, gorilla, echo, fiber")
	}
	return nil
}
