// Package commands implements CLI subcommands for the bluesky-collective tool.
package commands

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/viper"

	"github.com/consensuscode/bluesky-collective/pkg/atproto"
	"github.com/consensuscode/bluesky-collective/pkg/bluesky"
	"github.com/consensuscode/bluesky-collective/pkg/consensus"
	"github.com/consensuscode/bluesky-collective/pkg/storage"
)

// dataDir returns the base directory for local data storage.
// Priority: --data-dir flag / config file data_dir / BLUESKY_COLLECTIVE_DATA_DIR env var,
// then a relative default of "./.bluesky-collective-data" in the working directory.
func dataDir() string {
	dir := viper.GetString("data_dir")
	if dir != "" {
		return dir
	}

	return ".bluesky-collective-data"
}

// newConsensusChecker creates a FileChecker with configuration from viper.
func newConsensusChecker() (*consensus.FileChecker, error) {
	minPart := viper.GetInt("consensus.min_participation")
	if minPart < 1 {
		minPart = 3
	}

	timeout := viper.GetDuration("consensus.timeout")
	if timeout <= 0 {
		timeout = 24 * time.Hour
	}

	rules := consensus.NewDefaultRules(minPart, timeout)
	return consensus.NewFileChecker(filepath.Join(dataDir(), "consensus"), rules)
}

// newFileStore creates a FileStore with configuration from viper.
func newFileStore() (*storage.FileStore, error) {
	return storage.NewFileStore(filepath.Join(dataDir(), "storage"))
}

// newCollectiveClient creates a fully wired CollectiveClient.
func newCollectiveClient() (*bluesky.CollectiveClient, error) {
	checker, err := newConsensusChecker()
	if err != nil {
		return nil, fmt.Errorf("create consensus checker: %w", err)
	}

	store, err := newFileStore()
	if err != nil {
		return nil, fmt.Errorf("create storage: %w", err)
	}

	serviceURL := viper.GetString("bluesky.service")
	if serviceURL == "" {
		serviceURL = "https://bsky.social"
	}

	agentID := viper.GetString("agent.id")
	if agentID == "" {
		agentID = "unknown-agent"
	}

	atpClient := atproto.NewClient(serviceURL)
	poster := bluesky.NewATPAdapter(atpClient)

	return bluesky.NewCollectiveClient(poster, checker, store, agentID), nil
}

// currentAgentID returns the agent ID from configuration.
func currentAgentID() string {
	id := viper.GetString("agent.id")
	if id == "" {
		return "unknown-agent"
	}
	return id
}
