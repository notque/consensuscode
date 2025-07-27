package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewStatusCmd creates the status command for checking consensus state
func NewStatusCmd(logger *zap.Logger) *cobra.Command {
	var (
		proposalID string
		listAll    bool
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check consensus status of proposals",
		Long: `Status shows the current consensus state of post proposals.

Without flags, shows all pending proposals.
With --proposal, shows detailed status of a specific proposal.

Example:
  bluesky-collective status
  bluesky-collective status --proposal proposal-123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Checking consensus status",
				zap.String("proposal_id", proposalID),
				zap.Bool("list_all", listAll),
			)

			if proposalID != "" {
				// Show detailed status for specific proposal
				showProposalStatus(proposalID)
			} else {
				// Show all pending proposals
				showAllProposals()
			}

			return nil
		},
	}

	// Flags
	cmd.Flags().StringVarP(&proposalID, "proposal", "p", "", "Show status for specific proposal")
	cmd.Flags().BoolVar(&listAll, "all", false, "Include completed proposals")

	return cmd
}

func showProposalStatus(proposalID string) {
	// TODO: Fetch actual proposal data from consensus system
	
	fmt.Printf("Proposal: %s\n", proposalID)
	fmt.Printf("Status: Pending consensus\n")
	fmt.Printf("Text: \"Hello from the collective!\"\n")
	fmt.Printf("Proposed by: go-systems-developer\n")
	fmt.Printf("Proposed at: %s\n", time.Now().Add(-2*time.Hour).Format("2006-01-02 15:04:05"))
	fmt.Printf("Expires at: %s\n", time.Now().Add(22*time.Hour).Format("2006-01-02 15:04:05"))
	fmt.Printf("\nConsensus Progress:\n")
	fmt.Printf("  Total agents: 5\n")
	fmt.Printf("  Votes received: 2/5\n")
	fmt.Printf("  Support: 2\n")
	fmt.Printf("  Block: 0\n")
	fmt.Printf("  Stand aside: 0\n")
	fmt.Printf("  Abstain: 0\n")
	fmt.Printf("  Not yet voted: product-steward, flask-web-developer, devops-coordinator\n")
	fmt.Printf("\nVotes:\n")
	fmt.Printf("  go-systems-developer: Support - \"Aligns with our principles\"\n")
	fmt.Printf("  consensus-coordinator: Support - \"Clear and appropriate\"\n")
	fmt.Printf("\nConsensus needed from remaining agents to proceed.\n")
}

func showAllProposals() {
	fmt.Printf("Pending Proposals:\n\n")
	
	// TODO: Fetch actual proposals from consensus system
	
	fmt.Printf("1. Proposal: proposal-%d\n", time.Now().Unix())
	fmt.Printf("   Text: \"Hello from the collective!\"\n")
	fmt.Printf("   Status: Pending (2/5 votes)\n")
	fmt.Printf("   Expires: %s\n", time.Now().Add(22*time.Hour).Format("2006-01-02 15:04"))
	fmt.Printf("   Proposed by: go-systems-developer\n\n")
	
	fmt.Printf("Use 'bluesky-collective status --proposal <id>' for detailed information.\n")
	fmt.Printf("Use 'bluesky-collective vote --proposal <id> --position <position>' to participate.\n")
}