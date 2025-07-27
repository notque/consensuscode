package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"collectiveflow/internal/proposal"
)

// newProposalCmd creates the proposal command group
func newProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Manage collective proposals",
		Long: `Create, list, view, and manage proposals for collective consideration.

All proposals are created with collective consideration in mind - they require
systematic consultation and consensus building rather than individual approval.`,
	}

	// Add subcommands
	cmd.AddCommand(newProposalCreateCmd())
	cmd.AddCommand(newProposalListCmd())
	cmd.AddCommand(newProposalShowCmd())
	cmd.AddCommand(newProposalUpdateCmd())

	return cmd
}

// newProposalCreateCmd creates a new proposal
func newProposalCreateCmd() *cobra.Command {
	var (
		title       string
		description string
		urgency     string
		affectedAreas []string
	)

	cmd := &cobra.Command{
		Use:   "create [title]",
		Short: "Create a new proposal for collective consideration",
		Long: `Create a new proposal that will be submitted to the collective for
systematic consultation and consensus building.

The proposal will be assigned a unique ID and made available for all agents
to review and provide input on.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use argument title if provided, otherwise use flag
			if len(args) > 0 {
				title = args[0]
			}
			
			if title == "" {
				return fmt.Errorf("proposal title is required")
			}

			// Create the proposal
			prop := proposal.New{
				Title:         title,
				Description:   description,
				Urgency:       urgency,
				AffectedAreas: affectedAreas,
				Proposer:      "cli-user", // TODO: Get from config/context
				Date:          time.Now(),
			}

			result, err := proposal.Create(prop)
			if err != nil {
				return fmt.Errorf("failed to create proposal: %w", err)
			}

			fmt.Printf("✓ Proposal created successfully\n")
			fmt.Printf("  ID: %s\n", result.ID)
			fmt.Printf("  Title: %s\n", result.Title)
			fmt.Printf("  Status: %s\n", result.Status)
			fmt.Printf("  File: %s\n", result.FilePath)
			fmt.Printf("\nNext steps:\n")
			fmt.Printf("  1. Use 'collectiveflow consensus start %s' to begin consultation\n", result.ID)
			fmt.Printf("  2. Track progress with 'collectiveflow consensus status %s'\n", result.ID)

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&title, "title", "t", "", "Proposal title (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Detailed proposal description")
	cmd.Flags().StringVarP(&urgency, "urgency", "u", "medium", "Urgency level (low, medium, high, emergency)")
	cmd.Flags().StringSliceVarP(&affectedAreas, "affected", "a", []string{}, "Areas affected by this proposal")

	return cmd
}

// newProposalListCmd lists proposals
func newProposalListCmd() *cobra.Command {
	var (
		status     string
		limit      int
		showAll    bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		Long: `List proposals with optional filtering by status, urgency, or other criteria.

By default, shows active proposals that require attention.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := proposal.ListFilter{
				Status:  status,
				Limit:   limit,
				ShowAll: showAll,
			}

			proposals, err := proposal.List(filter)
			if err != nil {
				return fmt.Errorf("failed to list proposals: %w", err)
			}

			if len(proposals) == 0 {
				fmt.Printf("No proposals found matching criteria\n")
				return nil
			}

			fmt.Printf("Proposals (%d total):\n\n", len(proposals))
			for _, p := range proposals {
				fmt.Printf("  %s - %s\n", p.ID, p.Title)
				fmt.Printf("    Status: %s | Urgency: %s | Date: %s\n", 
					p.Status, p.Urgency, p.Date.Format("2006-01-02"))
				if len(p.AffectedAreas) > 0 {
					fmt.Printf("    Affects: %v\n", p.AffectedAreas)
				}
				fmt.Printf("\n")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&status, "status", "s", "", "Filter by status (proposed, consultation, consensus, implemented)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of proposals to show")
	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all proposals including completed ones")

	return cmd
}

// newProposalShowCmd shows a specific proposal
func newProposalShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [proposal-id]",
		Short: "Show detailed information about a proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			prop, err := proposal.Get(proposalID)
			if err != nil {
				return fmt.Errorf("failed to get proposal %s: %w", proposalID, err)
			}

			fmt.Printf("Proposal: %s\n", prop.Title)
			fmt.Printf("ID: %s\n", prop.ID)
			fmt.Printf("Status: %s\n", prop.Status)
			fmt.Printf("Urgency: %s\n", prop.Urgency)
			fmt.Printf("Proposer: %s\n", prop.Proposer)
			fmt.Printf("Date: %s\n", prop.Date.Format("2006-01-02 15:04:05"))
			
			if len(prop.AffectedAreas) > 0 {
				fmt.Printf("Affected Areas: %v\n", prop.AffectedAreas)
			}
			
			if prop.Description != "" {
				fmt.Printf("\nDescription:\n%s\n", prop.Description)
			}

			if prop.ConsensusStatus != "" {
				fmt.Printf("\nConsensus Status: %s\n", prop.ConsensusStatus)
			}

			return nil
		},
	}

	return cmd
}

// newProposalUpdateCmd updates a proposal
func newProposalUpdateCmd() *cobra.Command {
	var (
		title       string
		description string
		urgency     string
	)

	cmd := &cobra.Command{
		Use:   "update [proposal-id]",
		Short: "Update an existing proposal",
		Long: `Update the title, description, or urgency of an existing proposal.

Note: Significant changes may require restarting the consensus process.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			updates := proposal.Update{
				Title:       title,
				Description: description,
				Urgency:     urgency,
			}

			err := proposal.UpdateProposal(proposalID, updates)
			if err != nil {
				return fmt.Errorf("failed to update proposal %s: %w", proposalID, err)
			}

			fmt.Printf("✓ Proposal %s updated successfully\n", proposalID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Update proposal title")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Update proposal description")  
	cmd.Flags().StringVarP(&urgency, "urgency", "u", "", "Update urgency level")

	return cmd
}