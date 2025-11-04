package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"collectiveflow/internal/storage"
	"collectiveflow/internal/web"
)

// newWebCmd creates the web command group
func newWebCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Web interface commands",
		Long: `Manage the CollectiveFlow web interface.

The web interface provides a horizontal, accessible way to view proposals
and collective activity through a browser. No admin panels, no special
privileges - just transparent information access for all.`,
	}

	cmd.AddCommand(newWebServeCmd())

	return cmd
}

// newWebServeCmd creates the web serve command
func newWebServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web server",
		Long: `Start the embedded web server to access CollectiveFlow through a browser.

The web interface embodies horizontal principles:
  • No admin panels or privileged access
  • All proposals visible to everyone
  • Transparent consensus process visualization
  • Community bulletin board design
  • Mobile-responsive and accessible

The web server runs locally and serves the interface at http://localhost:8080
by default. All data is read from the same storage backend used by the CLI,
ensuring CLI and web interfaces have equal access and capabilities.`,
		RunE: runWebServe,
	}

	cmd.Flags().StringP("addr", "a", ":8080", "Address to bind the web server (e.g., ':8080', 'localhost:3000')")
	cmd.Flags().StringP("storage-path", "s", "", "Path to proposal storage (overrides config)")

	return cmd
}

// runWebServe executes the web serve command
func runWebServe(cmd *cobra.Command, args []string) error {
	// Get configuration
	addr, _ := cmd.Flags().GetString("addr")
	storagePath, _ := cmd.Flags().GetString("storage-path")

	// Use storage path from flag or config
	if storagePath == "" {
		storagePath = viper.GetString("storage.path")
	}

	// Ensure storage path exists
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return fmt.Errorf("storage path does not exist: %s\nCreate it with: mkdir -p %s", storagePath, storagePath)
	}

	// Initialize storage
	store, err := storage.NewFileStore(storagePath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Create web server
	server, err := web.NewServer(web.Config{
		Store: store,
		Addr:  addr,
	})
	if err != nil {
		return fmt.Errorf("failed to create web server: %w", err)
	}

	// Start server (blocks until stopped)
	if err := server.Start(); err != nil {
		return fmt.Errorf("web server error: %w", err)
	}

	return nil
}
