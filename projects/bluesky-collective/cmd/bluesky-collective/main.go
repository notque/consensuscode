// bluesky-collective is a CLI tool for managing the collective's Bluesky presence
// with mandatory consensus requirements for all posts.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/consensuscode/bluesky-collective/cmd/bluesky-collective/commands"
)

var (
	cfgFile string
	logger  *zap.Logger
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "bluesky-collective",
	Short: "Collective consensus-based Bluesky client",
	Long: `bluesky-collective is a tool for managing the collective's Bluesky presence
with mandatory consensus requirements for all posts.

All posts require consensus from the collective before being published.
This ensures horizontal accountability and collective decision-making.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bluesky-collective.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	
	// Bind flags to viper
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	// Initialize logger
	var err error
	logger, err = initLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewProposeCmd(logger))
	rootCmd.AddCommand(commands.NewVoteCmd(logger))
	rootCmd.AddCommand(commands.NewStatusCmd(logger))
	rootCmd.AddCommand(commands.NewPublishCmd(logger))
	rootCmd.AddCommand(commands.NewConfigCmd(logger))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bluesky-collective")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("BLUESKY_COLLECTIVE")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	
	// Set log level from config
	level := viper.GetString("log-level")
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", zap.Error(err))
		os.Exit(1)
	}
}