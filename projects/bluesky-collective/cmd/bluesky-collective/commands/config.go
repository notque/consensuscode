package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewConfigCmd creates the config command for managing settings
func NewConfigCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
		Long: `Config command manages the collective's Bluesky configuration.

Settings include:
- Bluesky credentials and endpoints
- Consensus timeouts and rules
- Agent identification
- Storage backends`,
	}

	cmd.AddCommand(newConfigShowCmd(logger))
	cmd.AddCommand(newConfigSetCmd(logger))
	cmd.AddCommand(newConfigInitCmd(logger))

	return cmd
}

func newConfigShowCmd(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Displaying configuration")

			fmt.Printf("Configuration File: %s\n", viper.ConfigFileUsed())
			fmt.Printf("\nBluesky Settings:\n")
			fmt.Printf("  Handle: %s\n", viper.GetString("bluesky.handle"))
			fmt.Printf("  Service: %s\n", viper.GetString("bluesky.service"))
			fmt.Printf("  Authenticated: %v\n", viper.GetString("bluesky.session") != "")
			
			fmt.Printf("\nConsensus Settings:\n")
			fmt.Printf("  Timeout: %s\n", viper.GetDuration("consensus.timeout"))
			fmt.Printf("  Minimum participation: %d\n", viper.GetInt("consensus.min_participation"))
			fmt.Printf("  Storage backend: %s\n", viper.GetString("consensus.storage"))
			
			fmt.Printf("\nAgent Settings:\n")
			fmt.Printf("  Agent ID: %s\n", viper.GetString("agent.id"))
			fmt.Printf("  Log level: %s\n", viper.GetString("log-level"))

			return nil
		},
	}
}

func newConfigSetCmd(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			logger.Info("Setting configuration value",
				zap.String("key", key),
				zap.String("value", value),
			)

			viper.Set(key, value)
			
			// Write config file
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("Set %s = %s\n", key, value)
			return nil
		},
	}
}

func newConfigInitCmd(logger *zap.Logger) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Initializing configuration")

			// Create default config
			viper.SetDefault("bluesky.service", "https://bsky.social")
			viper.SetDefault("bluesky.handle", "collectiveflow.bsky.social")
			viper.SetDefault("consensus.timeout", "24h")
			viper.SetDefault("consensus.min_participation", 3)
			viper.SetDefault("consensus.storage", "file")
			viper.SetDefault("agent.id", "go-systems-developer")
			viper.SetDefault("log-level", "info")

			// Determine config file location: use --path flag, then --config
			// root flag, then fall back to .bluesky-collective.yaml in the
			// current directory (no hardcoded home-directory path).
			if configPath == "" {
				if existing := viper.ConfigFileUsed(); existing != "" {
					configPath = existing
				} else {
					configPath = ".bluesky-collective.yaml"
				}
			}

			// Create parent directory if needed
			dir := filepath.Dir(configPath)
			if dir != "." && dir != "" {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create config directory: %w", err)
				}
			}

			viper.SetConfigFile(configPath)

			// Write config file
			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			fmt.Printf("Configuration initialized at: %s\n", configPath)
			fmt.Printf("\nNext steps:\n")
			fmt.Printf("1. Set your Bluesky credentials:\n")
			fmt.Printf("   bluesky-collective config set bluesky.identifier <your-handle>\n")
			fmt.Printf("   bluesky-collective config set bluesky.password <your-app-password>\n")
			fmt.Printf("2. Verify settings:\n")
			fmt.Printf("   bluesky-collective config show\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "path", "p", "", "Path for the config file (default: ./.bluesky-collective.yaml)")

	return cmd
}