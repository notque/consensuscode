package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"collectiveflow/internal/proposal"
)

// newConsensusCmd creates the consensus command group
func newConsensusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consensus",
		Short: "Manage consensus processes for proposals",
		Long: `Start, track, and complete consensus processes for collective proposals.

The consensus process ensures all collective members have an opportunity to
provide input and that decisions are made through genuine collective agreement
rather than hierarchical authority.`,
	}

	// Add subcommands
	cmd.AddCommand(newConsensusStartCmd())
	cmd.AddCommand(newConsensusStatusCmd())
	cmd.AddCommand(newConsensusInputCmd())
	cmd.AddCommand(newConsensusCompleteCmd())

	return cmd
}

// newConsensusStartCmd starts a consensus process for a proposal
func newConsensusStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [proposal-id]",
		Short: "Start consensus process for a proposal",
		Long: `Begin the consensus-building process for a proposal.

This moves the proposal into consultation status and begins tracking
input from all collective members. Remember - consensus is not voting,
it's a process of addressing concerns until all can support the decision.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			// Update proposal status to consultation
			err := proposal.UpdateStatus(proposalID, proposal.StatusConsultation, "cli-user")
			if err != nil {
				return fmt.Errorf("failed to start consensus process: %w", err)
			}

			fmt.Printf("✓ Consensus process started for proposal %s\n", proposalID)
			fmt.Printf("\nNext steps:\n")
			fmt.Printf("  1. Notify all collective members to review the proposal\n")
			fmt.Printf("  2. Gather input using 'collectiveflow consensus input %s'\n", proposalID)
			fmt.Printf("  3. Address any concerns raised by members\n")
			fmt.Printf("  4. Complete consensus when all concerns are addressed\n")

			return nil
		},
	}

	return cmd
}

// newConsensusStatusCmd shows the status of a consensus process
func newConsensusStatusCmd() *cobra.Command {
	var showDetails bool

	cmd := &cobra.Command{
		Use:   "status [proposal-id]",
		Short: "Show consensus status for a proposal",
		Long: `Display the current status of the consensus process for a proposal,
including who has provided input and any outstanding concerns.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			// Get the proposal
			prop, err := proposal.Get(proposalID)
			if err != nil {
				return fmt.Errorf("failed to get proposal: %w", err)
			}

			fmt.Printf("Consensus Status for: %s\n", prop.Title)
			fmt.Printf("ID: %s\n", prop.ID)
			fmt.Printf("Status: %s\n", prop.Status)
			
			if prop.ConsensusStatus != "" {
				fmt.Printf("Consensus: %s\n", prop.ConsensusStatus)
			}

			// Show consultations summary
			if len(prop.Consultations) > 0 {
				fmt.Printf("\nConsultations received: %d\n", len(prop.Consultations))
				
				supportCount := 0
				var concerns []string
				
				for _, consultation := range prop.Consultations {
					if consultation.Support {
						supportCount++
					} else {
						concerns = append(concerns, consultation.Concerns...)
					}
				}
				
				fmt.Printf("  Supporting: %d\n", supportCount)
				fmt.Printf("  With concerns: %d\n", len(prop.Consultations)-supportCount)
				
				if len(concerns) > 0 {
					fmt.Printf("\nOutstanding concerns:\n")
					for i, concern := range concerns {
						fmt.Printf("  %d. %s\n", i+1, concern)
					}
				}
				
				if showDetails {
					fmt.Printf("\nDetailed consultations:\n")
					for _, c := range prop.Consultations {
						fmt.Printf("\n  From: %s (at %s)\n", c.Contributor, c.Timestamp.Format("2006-01-02 15:04"))
						fmt.Printf("  Support: %v\n", c.Support)
						fmt.Printf("  Input: %s\n", c.Input)
						if len(c.Concerns) > 0 {
							fmt.Printf("  Concerns: %v\n", c.Concerns)
						}
					}
				}
			} else {
				fmt.Printf("\nNo consultations received yet.\n")
			}

			// Check if consensus is possible
			if prop.Status == proposal.StatusConsultation {
				if prop.HasUnanimousSupport() {
					fmt.Printf("\n✓ All consultations support the proposal - consensus possible!\n")
				} else {
					blockingConcerns := prop.GetBlockingConcerns()
					if len(blockingConcerns) > 0 {
						fmt.Printf("\n⚠ Consensus blocked by concerns - address these before proceeding.\n")
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&showDetails, "details", "d", false, "Show detailed consultation information")

	return cmd
}

// newConsensusInputCmd records consultation input
func newConsensusInputCmd() *cobra.Command {
	var (
		contributor string
		input       string
		support     bool
		concerns    []string
	)

	cmd := &cobra.Command{
		Use:   "input [proposal-id]",
		Short: "Provide consultation input for a proposal",
		Long: `Record your input during the consensus process.

This allows collective members to express support or concerns about a proposal.
Remember that consensus is about addressing concerns, not simply voting.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			if contributor == "" {
				return fmt.Errorf("contributor name is required")
			}

			if input == "" {
				return fmt.Errorf("input text is required")
			}

			// Create consultation record
			consultation := proposal.Consultation{
				Contributor: contributor,
				Timestamp:   time.Now(),
				Input:       input,
				Support:     support,
				Concerns:    concerns,
			}

			// Add the consultation
			err := proposal.AddConsultationInput(proposalID, consultation)
			if err != nil {
				return fmt.Errorf("failed to record consultation: %w", err)
			}

			fmt.Printf("✓ Consultation input recorded for proposal %s\n", proposalID)
			
			if support {
				fmt.Printf("  Your support has been noted.\n")
			} else {
				fmt.Printf("  Your concerns have been recorded for collective consideration.\n")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&contributor, "contributor", "c", "", "Name of the contributor (required)")
	cmd.Flags().StringVarP(&input, "input", "i", "", "Consultation input text (required)")
	cmd.Flags().BoolVarP(&support, "support", "s", false, "Indicate support for the proposal")
	cmd.Flags().StringSliceVar(&concerns, "concerns", []string{}, "List specific concerns (comma-separated)")

	cmd.MarkFlagRequired("contributor")
	cmd.MarkFlagRequired("input")

	return cmd
}

// newConsensusCompleteCmd completes a consensus process
func newConsensusCompleteCmd() *cobra.Command {
	var (
		result    string
		rationale string
	)

	cmd := &cobra.Command{
		Use:   "complete [proposal-id]",
		Short: "Complete the consensus process for a proposal",
		Long: `Record the collective's decision after the consensus process.

This should only be done after all members have had opportunity to provide
input and any concerns have been addressed. The decision reflects the
collective will, not individual authority.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			if result == "" {
				return fmt.Errorf("decision result is required (approved/rejected/deferred/no_consensus)")
			}

			if rationale == "" {
				return fmt.Errorf("rationale for the decision is required")
			}

			// Parse decision result
			var decisionResult proposal.DecisionResult
			switch result {
			case "approved":
				decisionResult = proposal.DecisionApproved
			case "rejected":
				decisionResult = proposal.DecisionRejected
			case "deferred":
				decisionResult = proposal.DecisionDeferred
			case "no_consensus":
				decisionResult = proposal.DecisionNoConsensus
			default:
				return fmt.Errorf("invalid decision result: %s (use approved/rejected/deferred/no_consensus)", result)
			}

			// Create decision record
			decision := proposal.Decision{
				Result:    decisionResult,
				Timestamp: time.Now(),
				Rationale: rationale,
			}

			// Record the decision
			err := proposal.RecordDecision(proposalID, decision)
			if err != nil {
				return fmt.Errorf("failed to record decision: %w", err)
			}

			fmt.Printf("✓ Consensus process completed for proposal %s\n", proposalID)
			fmt.Printf("  Decision: %s\n", result)
			fmt.Printf("  Rationale: %s\n", rationale)

			switch decisionResult {
			case proposal.DecisionApproved:
				fmt.Printf("\nThe proposal has been approved by collective consensus.\n")
				fmt.Printf("Next step: Implement the proposal collaboratively.\n")
			case proposal.DecisionRejected:
				fmt.Printf("\nThe proposal has been rejected by the collective.\n")
			case proposal.DecisionDeferred:
				fmt.Printf("\nThe proposal has been deferred for future consideration.\n")
			case proposal.DecisionNoConsensus:
				fmt.Printf("\nNo consensus could be reached at this time.\n")
				fmt.Printf("Consider revisiting after addressing outstanding concerns.\n")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&result, "result", "r", "", "Decision result (approved/rejected/deferred/no_consensus)")
	cmd.Flags().StringVar(&rationale, "rationale", "", "Rationale for the collective decision")

	cmd.MarkFlagRequired("result")
	cmd.MarkFlagRequired("rationale")

	return cmd
}