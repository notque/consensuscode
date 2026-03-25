package main

import (
	"os"

	"collectiveflow/internal/cli"
)

// Build-time variables set via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app := cli.NewApp(version, commit, date)
	if err := app.Execute(); err != nil {
		// Cobra already prints the error, so just exit with a non-zero code.
		os.Exit(1)
	}
}
