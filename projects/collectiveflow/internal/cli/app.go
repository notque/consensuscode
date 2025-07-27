package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// App represents the CollectiveFlow CLI application
type App struct {
	rootCmd *cobra.Command
	version string
	commit  string
	date    string
}

// NewApp creates a new CollectiveFlow CLI application
func NewApp(version, commit, date string) *App {
	app := &App{
		version: version,
		commit:  commit,
		date:    date,
	}

	app.rootCmd = &cobra.Command{
		Use:   "collectiveflow",
		Short: "Horizontal collective decision-making support tool",
		Long: `CollectiveFlow supports horizontal collective decision-making processes
through consensus coordination, proposal management, and decision documentation.

Built through genuine consensus by a libertarian socialist AI agent collective,
this tool demonstrates that high-quality software can be developed through
horizontal coordination without hierarchy.`,
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
	}

	// Initialize configuration
	app.initConfig()
	
	// Add command groups
	app.addCommands()

	return app
}

// Execute runs the CLI application
func (a *App) Execute() error {
	return a.rootCmd.Execute()
}

// initConfig initializes the application configuration
func (a *App) initConfig() {
	// Set configuration defaults that reflect collective values
	viper.SetDefault("collective.name", "consensus-collective")
	viper.SetDefault("storage.type", "file")
	viper.SetDefault("storage.path", "./collective-data")
	viper.SetDefault("consensus.require_all_agents", true)
	viper.SetDefault("proposals.auto_id", true)
	
	// Configuration file locations (collective can modify these)
	viper.SetConfigName("collectiveflow")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.collectiveflow")
	viper.AddConfigPath("/etc/collectiveflow")

	// Read configuration if available
	if err := viper.ReadInConfig(); err != nil {
		// Config file not required - defaults work
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Warning: Error reading config file: %v\n", err)
		}
	}

	// Allow environment variable overrides
	viper.SetEnvPrefix("COLLECTIVEFLOW")
	viper.AutomaticEnv()
}

// addCommands adds all command groups to the root command
func (a *App) addCommands() {
	// Add proposal management commands
	a.rootCmd.AddCommand(newProposalCmd())
	
	// Add consensus process commands
	a.rootCmd.AddCommand(newConsensusCmd())
	
	// Add collective status and information commands
	a.rootCmd.AddCommand(newStatusCmd())
	
	// Add configuration management commands
	a.rootCmd.AddCommand(newConfigCmd())
}