package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewStatusCmd creates the status command for checking consensus state.
func NewStatusCmd(logger *zap.Logger) *cobra.Command {
	var (
		proposalID string
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
			)

			checker, err := newConsensusChecker()
			if err != nil {
				return fmt.Errorf("initialize consensus checker: %w", err)
			}

			if proposalID != "" {
				// Show detailed status for specific proposal
				proposal, err := checker.GetProposal(cmd.Context(), proposalID)
				if err != nil {
					return fmt.Errorf("get proposal: %w", err)
				}

				decision, err := checker.GetDecision(cmd.Context(), proposalID)
				if err != nil {
					return fmt.Errorf("get decision: %w", err)
				}

				fmt.Printf("Proposal: %s\n", proposal.ID)
				fmt.Printf("  Status: %s\n", decision.Status)
				fmt.Printf("  Text: %q\n", proposal.Content)
				fmt.Printf("  Proposed by: %s\n", proposal.ProposedBy)
				fmt.Printf("  Reasoning: %s\n", proposal.Reasoning)
				fmt.Printf("  Proposed at: %s\n", proposal.ProposedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("  Expires at: %s\n", proposal.ExpiresAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("\nVotes (%d):\n", len(decision.AgentVotes))

				if len(decision.AgentVotes) == 0 {
					fmt.Printf("  (no votes yet)\n")
				}
				for agentID, vote := range decision.AgentVotes {
					fmt.Printf("  %s: %s", agentID, vote.Position)
					if vote.Reasoning != "" {
						fmt.Printf(" - %q", vote.Reasoning)
					}
					fmt.Println()
				}

				if decision.ConsensusAt != nil {
					fmt.Printf("\nConsensus reached at: %s\n", decision.ConsensusAt.Format("2006-01-02 15:04:05"))
				}
			} else {
				// Show all pending proposals
				proposals, err := checker.ListPendingProposals(cmd.Context())
				if err != nil {
					return fmt.Errorf("list proposals: %w", err)
				}

				if len(proposals) == 0 {
					fmt.Println("No pending proposals.")
					fmt.Println("\nUse 'bluesky-collective propose' to create one.")
					return nil
				}

				fmt.Printf("Pending Proposals (%d):\n\n", len(proposals))
				for i, proposal := range proposals {
					decision, err := checker.GetDecision(cmd.Context(), proposal.ID)
					voteCount := 0
					status := "unknown"
					if err == nil {
						voteCount = len(decision.AgentVotes)
						status = string(decision.Status)
					}

					fmt.Printf("%d. %s\n", i+1, proposal.ID)
					fmt.Printf("   Text: %q\n", proposal.Content)
					fmt.Printf("   Status: %s (%d votes)\n", status, voteCount)
					fmt.Printf("   Proposed by: %s\n", proposal.ProposedBy)
					fmt.Printf("   Expires: %s\n\n", proposal.ExpiresAt.Format("2006-01-02 15:04"))
				}

				fmt.Println("Use 'bluesky-collective status --proposal <id>' for details.")
				fmt.Println("Use 'bluesky-collective vote --proposal <id> --position <position>' to participate.")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&proposalID, "proposal", "p", "", "Show status for specific proposal")

	return cmd
}
