package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newConfigCmd creates the config command group
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CollectiveFlow configuration",
		Long: `View and manage CollectiveFlow configuration settings.

Configuration can be modified through collective consensus. All changes
should be proposed and agreed upon by the collective rather than imposed
by individual preference.`,
	}

	// Add subcommands
	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigInitCmd())
	cmd.AddCommand(newConfigPathCmd())

	return cmd
}

// newConfigShowCmd shows current configuration
func newConfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long: `Display all current configuration values, including defaults
and any overrides from config files or environment variables.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("CollectiveFlow Configuration")
			fmt.Println("============================")

			// Show config file location if found
			configFile := viper.ConfigFileUsed()
			if configFile != "" {
				fmt.Printf("Config file: %s\n\n", configFile)
			} else {
				fmt.Println("No config file found - using defaults\n")
			}

			// Display all settings
			fmt.Println("Collective Settings:")
			fmt.Printf("  collective.name: %s\n", viper.GetString("collective.name"))
			fmt.Printf("  consensus.require_all_agents: %v\n", viper.GetBool("consensus.require_all_agents"))

			fmt.Println("\nStorage Settings:")
			fmt.Printf("  storage.type: %s\n", viper.GetString("storage.type"))
			fmt.Printf("  storage.path: %s\n", viper.GetString("storage.path"))

			fmt.Println("\nProposal Settings:")
			fmt.Printf("  proposals.auto_id: %v\n", viper.GetBool("proposals.auto_id"))

			// Show environment variable overrides
			fmt.Println("\nEnvironment Variable Overrides:")
			envVars := []string{
				"COLLECTIVEFLOW_COLLECTIVE_NAME",
				"COLLECTIVEFLOW_STORAGE_TYPE",
				"COLLECTIVEFLOW_STORAGE_PATH",
			}

			foundEnv := false
			for _, env := range envVars {
				if val := os.Getenv(env); val != "" {
					fmt.Printf("  %s=%s\n", env, val)
					foundEnv = true
				}
			}

			if !foundEnv {
				fmt.Println("  (none)")
			}

			fmt.Println("\nNote: Configuration changes should be made through collective consensus.")

			return nil
		},
	}

	return cmd
}

// newConfigInitCmd initializes a configuration file
func newConfigInitCmd() *cobra.Command {
	var (
		path      string
		force     bool
		storePath string
		name      string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a configuration file",
		Long: `Create a new configuration file with default values.

This file can be modified by the collective to customize behavior.
Remember that configuration changes affect all collective members
and should be made through consensus.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine config file path
			if path == "" {
				path = "./collectiveflow.yaml"
			}

			// Check if file exists
			if _, err := os.Stat(path); err == nil && !force {
				return fmt.Errorf("config file already exists at %s (use --force to overwrite)", path)
			}


			// Create directory if needed
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			// Write config file
			file, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
			defer file.Close()

			// Write YAML with custom formatting
			fmt.Fprintln(file, "# CollectiveFlow Configuration")
			fmt.Fprintln(file, "# This file was created through collective consensus")
			fmt.Fprintln(file, "# Modify only through collective agreement")
			fmt.Fprintln(file)
			fmt.Fprintln(file, "collective:")
			fmt.Fprintf(file, "  name: %s\n", name)
			fmt.Fprintln(file, "  # Name of your collective")
			fmt.Fprintln(file)
			fmt.Fprintln(file, "storage:")
			fmt.Fprintln(file, "  type: file")
			fmt.Fprintf(file, "  path: %s\n", storePath)
			fmt.Fprintln(file, "  # Storage backend configuration")
			fmt.Fprintln(file)
			fmt.Fprintln(file, "consensus:")
			fmt.Fprintln(file, "  require_all_agents: true")
			fmt.Fprintln(file, "  # Whether all agents must be consulted")
			fmt.Fprintln(file)
			fmt.Fprintln(file, "proposals:")
			fmt.Fprintln(file, "  auto_id: true")
			fmt.Fprintln(file, "  # Automatically generate proposal IDs")

			fmt.Printf("✓ Configuration file created at %s\n", path)
			fmt.Println("\nNext steps:")
			fmt.Println("  1. Review the configuration with your collective")
			fmt.Println("  2. Propose any changes through the consensus process")
			fmt.Println("  3. Restart CollectiveFlow to use the new configuration")

			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Path for the config file (default: ./collectiveflow.yaml)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing config file")
	cmd.Flags().StringVar(&storePath, "storage-path", "./collective-data", "Path for data storage")
	cmd.Flags().StringVar(&name, "collective-name", "consensus-collective", "Name of your collective")

	return cmd
}

// newConfigPathCmd shows configuration file search paths
func newConfigPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Show configuration file search paths",
		Long: `Display all paths where CollectiveFlow searches for configuration files.

Files are searched in order, with the first found file taking precedence.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Configuration File Search Paths:")
			fmt.Println("================================")

			// Get search paths
			paths := []string{
				"./collectiveflow.yaml",
				"$HOME/.collectiveflow/collectiveflow.yaml",
				"/etc/collectiveflow/collectiveflow.yaml",
			}

			// Expand and check each path
			for i, p := range paths {
				expanded := os.ExpandEnv(p)
				fmt.Printf("%d. %s", i+1, expanded)

				// Check if file exists
				if _, err := os.Stat(expanded); err == nil {
					fmt.Printf(" ✓ (exists)")
					if viper.ConfigFileUsed() == expanded {
						fmt.Printf(" [ACTIVE]")
					}
				} else {
					fmt.Printf(" ✗ (not found)")
				}
				fmt.Println()
			}

			// Show which file is actually being used
			if configFile := viper.ConfigFileUsed(); configFile != "" {
				fmt.Printf("\nCurrently using: %s\n", configFile)
			} else {
				fmt.Println("\nNo configuration file found - using defaults")
			}

			fmt.Println("\nEnvironment variables override file settings:")
			fmt.Println("  Prefix: COLLECTIVEFLOW_")
			fmt.Println("  Example: COLLECTIVEFLOW_STORAGE_PATH=/tmp/collective")

			return nil
		},
	}

	return cmd
}