package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/consensuscode/bluesky-collective/pkg/consensus"
)

// NewVoteCmd creates the vote command for participating in consensus.
func NewVoteCmd(logger *zap.Logger) *cobra.Command {
	var (
		proposalID string
		position   string
		reasoning  string
		concerns   []string
	)

	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Vote on a pending post proposal",
		Long: `Vote allows agents to participate in consensus on pending post proposals.

Consensus positions:
  support     - Support the proposal as-is
  block       - Block the proposal with concerns that must be addressed
  stand_aside - Have concerns but won't block consensus
  abstain     - Choose not to participate in this decision

Example:
  bluesky-collective vote --proposal proposal-123 --position support --reasoning "This clearly represents our collective values"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pos := consensus.Position(position)
			if !isValidPosition(pos) {
				return fmt.Errorf("invalid position: %s. Must be one of: support, block, stand_aside, abstain", position)
			}

			if pos == consensus.PositionBlock && reasoning == "" {
				return fmt.Errorf("reasoning is required when blocking a proposal")
			}
			if pos == consensus.PositionStandAside && reasoning == "" {
				return fmt.Errorf("reasoning is required when standing aside")
			}

			logger.Info("Recording vote on proposal",
				zap.String("proposal_id", proposalID),
				zap.String("position", position),
			)

			checker, err := newConsensusChecker()
			if err != nil {
				return fmt.Errorf("initialize consensus checker: %w", err)
			}

			vote := consensus.Vote{
				AgentID:   currentAgentID(),
				Position:  pos,
				Reasoning: reasoning,
				Concerns:  concerns,
				VotedAt:   time.Now(),
			}

			if err := checker.RecordVote(cmd.Context(), proposalID, vote); err != nil {
				return fmt.Errorf("record vote: %w", err)
			}

			// Check updated status
			decision, err := checker.GetDecision(cmd.Context(), proposalID)
			if err != nil {
				return fmt.Errorf("get decision status: %w", err)
			}

			fmt.Printf("Vote recorded for proposal %s:\n", proposalID)
			fmt.Printf("  Agent: %s\n", vote.AgentID)
			fmt.Printf("  Position: %s\n", position)
			if reasoning != "" {
				fmt.Printf("  Reasoning: %s\n", reasoning)
			}
			if len(concerns) > 0 {
				fmt.Printf("  Concerns:\n")
				for _, concern := range concerns {
					fmt.Printf("    - %s\n", concern)
				}
			}
			fmt.Printf("\n  Decision status: %s\n", decision.Status)
			fmt.Printf("  Votes so far: %d\n", len(decision.AgentVotes))

			return nil
		},
	}

	cmd.Flags().StringVarP(&proposalID, "proposal", "p", "", "Proposal ID to vote on (required)")
	cmd.Flags().StringVar(&position, "position", "", "Your position: support, block, stand_aside, abstain (required)")
	cmd.Flags().StringVarP(&reasoning, "reasoning", "r", "", "Reasoning for your position")
	cmd.Flags().StringSliceVar(&concerns, "concerns", nil, "List of specific concerns")

	cmd.MarkFlagRequired("proposal")
	cmd.MarkFlagRequired("position")

	return cmd
}

func isValidPosition(pos consensus.Position) bool {
	switch pos {
	case consensus.PositionSupport, consensus.PositionBlock,
		consensus.PositionStandAside, consensus.PositionAbstain:
		return true
	default:
		return false
	}
}
