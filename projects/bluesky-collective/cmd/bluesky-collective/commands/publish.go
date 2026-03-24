package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewPublishCmd creates the publish command for posting after consensus.
func NewPublishCmd(logger *zap.Logger) *cobra.Command {
	var (
		proposalID string
		dryRun     bool
	)

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish a post after consensus is reached",
		Long: `Publish posts a proposal to Bluesky after consensus has been reached.

This command will only succeed if:
1. Consensus has been reached on the proposal
2. No agents have blocked the proposal
3. Valid Bluesky credentials are configured

Example:
  bluesky-collective publish --proposal proposal-123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Attempting to publish proposal",
				zap.String("proposal_id", proposalID),
				zap.Bool("dry_run", dryRun),
			)

			// Check consensus first using the checker directly
			checker, err := newConsensusChecker()
			if err != nil {
				return fmt.Errorf("initialize consensus checker: %w", err)
			}

			decision, err := checker.GetDecision(cmd.Context(), proposalID)
			if err != nil {
				return fmt.Errorf("get decision: %w", err)
			}

			fmt.Printf("Proposal: %s\n", proposalID)
			fmt.Printf("Consensus status: %s\n", decision.Status)
			fmt.Printf("Votes: %d\n", len(decision.AgentVotes))

			if decision.Status != "consensus" {
				fmt.Printf("\nCannot publish: consensus not yet reached.\n")
				fmt.Printf("Use 'bluesky-collective status --proposal %s' to check progress.\n", proposalID)
				return fmt.Errorf("consensus not reached (status: %s)", decision.Status)
			}

			if dryRun {
				store, err := newFileStore()
				if err != nil {
					return fmt.Errorf("initialize storage: %w", err)
				}
				req, err := store.GetPostRequest(cmd.Context(), proposalID)
				if err != nil {
					return fmt.Errorf("get post request: %w", err)
				}
				fmt.Printf("\n[DRY RUN] Would publish:\n")
				fmt.Printf("  Text: %q\n", req.Text)
				fmt.Printf("  Consensus ID: %s\n", proposalID)
				return nil
			}

			// Check credentials
			identifier := viper.GetString("bluesky.identifier")
			password := viper.GetString("bluesky.password")
			if identifier == "" || password == "" {
				return fmt.Errorf("Bluesky credentials not configured. Use 'bluesky-collective config set bluesky.identifier <handle>' and 'bluesky-collective config set bluesky.password <app-password>'")
			}

			client, err := newCollectiveClient()
			if err != nil {
				return fmt.Errorf("initialize client: %w", err)
			}

			// Authenticate
			if err := client.Authenticate(cmd.Context(), identifier, password); err != nil {
				return fmt.Errorf("authenticate with Bluesky: %w", err)
			}

			// Publish
			result, err := client.PublishWithConsensus(cmd.Context(), proposalID)
			if err != nil {
				return fmt.Errorf("publish: %w", err)
			}

			fmt.Printf("\nPost published successfully!\n")
			fmt.Printf("  URI: %s\n", result.URI)
			fmt.Printf("  CID: %s\n", result.CID)
			fmt.Printf("  Consensus ID: %s\n", result.ConsensusID)

			return nil
		},
	}

	cmd.Flags().StringVarP(&proposalID, "proposal", "p", "", "Proposal ID to publish (required)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be published without actually posting")

	cmd.MarkFlagRequired("proposal")

	return cmd
}
