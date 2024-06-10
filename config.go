package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	DockerRunning bool   `mapstructure:"DOCKER_RUNNING"`
	DataFolder    string `mapstructure:"DATA_FOLDER"`
	UseFsNotify   bool   `mapstructure:"USE_FSNOTIFY"`
	UsePolling    bool   `mapstructure:"USE_POLLING"`
	PollingTime   int    `mapstructure:"POLLING_TIME"`
}

// Get the configuration from order:
// - .env file
// - the environment variables
// - the command line flags
func getEnvConfig() Config {
	// Initialize Viper
	viper.SetConfigFile(".env") // Load .env file
	viper.AutomaticEnv()        // Read environment variables

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

	// Bind flags to Viper keys
	viper.BindPFlag("DATA_FOLDER", rootCmd.Flags().Lookup("data-folder"))
	viper.BindPFlag("USE_FSNOTIFY", rootCmd.Flags().Lookup("use-fsnotify"))
	viper.BindPFlag("USE_POLLING", rootCmd.Flags().Lookup("use-polling"))
	viper.BindPFlag("POLLING_TIME", rootCmd.Flags().Lookup("polling-time"))

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
	}
}

func (config Config) Validate() error {
	if config.DataFolder == "" {
		return errors.New("data folder is required")
	}

	if config.UsePolling && config.PollingTime <= 0 {
		return errors.New("polling time should be greater than 0")
	}

	return nil
}
