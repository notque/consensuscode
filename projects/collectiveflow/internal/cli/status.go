package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"collectiveflow/internal/proposal"
)

// newStatusCmd creates the status command group
func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show collective status and information",
		Long: `Display information about the collective's current state,
including active proposals, recent decisions, and overall health.

This provides transparency into collective operations without
creating hierarchical oversight.`,
		RunE: showOverallStatus,
	}

	// Add subcommands
	cmd.AddCommand(newStatusActiveCmd())
	cmd.AddCommand(newStatusHistoryCmd())
	cmd.AddCommand(newStatusHealthCmd())

	return cmd
}

// showOverallStatus displays the overall collective status
func showOverallStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("CollectiveFlow Status")
	fmt.Println("====================")
	fmt.Printf("Version: %s\n", cmd.Root().Version)
	fmt.Printf("Date: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Get all proposals for statistics
	allProposals, err := proposal.List(proposal.ListFilter{ShowAll: true})
	if err != nil {
		return fmt.Errorf("failed to get proposals: %w", err)
	}

	// Count by status
	statusCounts := make(map[proposal.ProposalStatus]int)
	for _, p := range allProposals {
		statusCounts[p.Status]++
	}

	fmt.Println("Proposal Statistics:")
	fmt.Printf("  Total proposals: %d\n", len(allProposals))
	fmt.Printf("  Proposed: %d\n", statusCounts[proposal.StatusProposed])
	fmt.Printf("  In consultation: %d\n", statusCounts[proposal.StatusConsultation])
	fmt.Printf("  Consensus reached: %d\n", statusCounts[proposal.StatusConsensus])
	fmt.Printf("  Implemented: %d\n", statusCounts[proposal.StatusImplemented])
	fmt.Printf("  Withdrawn: %d\n", statusCounts[proposal.StatusWithdrawn])
	fmt.Printf("  Blocked: %d\n", statusCounts[proposal.StatusBlocked])

	// Show recent activity
	fmt.Println("\nRecent Activity:")
	recentCount := 0
	for _, p := range allProposals {
		if time.Since(p.Date) < 7*24*time.Hour {
			recentCount++
			fmt.Printf("  - %s: %s (%s)\n", p.ID, p.Title, p.Status)
			if recentCount >= 5 {
				break
			}
		}
	}

	if recentCount == 0 {
		fmt.Println("  No proposals in the last 7 days")
	}

	// Show proposals needing attention
	needingAttention := 0
	fmt.Println("\nProposals Needing Attention:")
	for _, p := range allProposals {
		if p.Status == proposal.StatusProposed || p.Status == proposal.StatusConsultation {
			needingAttention++
			urgencyIcon := ""
			switch p.Urgency {
			case proposal.UrgencyEmergency:
				urgencyIcon = "🚨"
			case proposal.UrgencyHigh:
				urgencyIcon = "⚠️ "
			case proposal.UrgencyMedium:
				urgencyIcon = "📌"
			case proposal.UrgencyLow:
				urgencyIcon = "💭"
			}
			fmt.Printf("  %s %s: %s (%s)\n", urgencyIcon, p.ID, p.Title, p.Status)
		}
	}

	if needingAttention == 0 {
		fmt.Println("  No proposals currently need attention")
	}

	fmt.Println("\nUse 'collectiveflow status active' for more details on active proposals")

	return nil
}

// newStatusActiveCmd shows active proposals
func newStatusActiveCmd() *cobra.Command {
	var (
		urgencyFilter string
		limit         int
	)

	cmd := &cobra.Command{
		Use:   "active",
		Short: "Show active proposals requiring collective attention",
		Long: `Display all proposals that are currently active and may require
collective input or action.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get active proposals
			filter := proposal.ListFilter{
				Urgency: urgencyFilter,
				Limit:   limit,
			}

			proposals, err := proposal.List(filter)
			if err != nil {
				return fmt.Errorf("failed to list active proposals: %w", err)
			}

			if len(proposals) == 0 {
				fmt.Println("No active proposals found.")
				return nil
			}

			fmt.Printf("Active Proposals (%d total):\n\n", len(proposals))

			// Group by status
			byStatus := make(map[proposal.ProposalStatus][]*proposal.Proposal)
			for _, p := range proposals {
				byStatus[p.Status] = append(byStatus[p.Status], p)
			}

			// Show consultation proposals first (need immediate attention)
			if consultation := byStatus[proposal.StatusConsultation]; len(consultation) > 0 {
				fmt.Println("In Consultation (needs input):")
				for _, p := range consultation {
					showProposalSummary(p)
				}
				fmt.Println()
			}

			// Then proposed
			if proposed := byStatus[proposal.StatusProposed]; len(proposed) > 0 {
				fmt.Println("Proposed (awaiting consensus process):")
				for _, p := range proposed {
					showProposalSummary(p)
				}
				fmt.Println()
			}

			// Then consensus reached
			if consensus := byStatus[proposal.StatusConsensus]; len(consensus) > 0 {
				fmt.Println("Consensus Reached (ready to implement):")
				for _, p := range consensus {
					showProposalSummary(p)
				}
				fmt.Println()
			}

			// Then blocked
			if blocked := byStatus[proposal.StatusBlocked]; len(blocked) > 0 {
				fmt.Println("Blocked (needs concern resolution):")
				for _, p := range blocked {
					showProposalSummary(p)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&urgencyFilter, "urgency", "u", "", "Filter by urgency (low/medium/high/emergency)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "Maximum number of proposals to show")

	return cmd
}

// newStatusHistoryCmd shows collective decision history
func newStatusHistoryCmd() *cobra.Command {
	var (
		days  int
		limit int
	)

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show collective decision history",
		Long: `Display the history of collective decisions and actions,
providing transparency into how the collective has operated.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all proposals
			allProposals, err := proposal.List(proposal.ListFilter{ShowAll: true})
			if err != nil {
				return fmt.Errorf("failed to get proposals: %w", err)
			}

			// Filter to completed proposals within time range
			cutoff := time.Now().AddDate(0, 0, -days)
			var completed []*proposal.Proposal

			for _, p := range allProposals {
				if p.Decision != nil && p.Decision.Timestamp.After(cutoff) {
					completed = append(completed, p)
				}
			}

			if len(completed) == 0 {
				fmt.Printf("No decisions recorded in the last %d days.\n", days)
				return nil
			}

			fmt.Printf("Collective Decision History (last %d days):\n\n", days)

			count := 0
			for _, p := range completed {
				if count >= limit {
					fmt.Printf("\n... and %d more decisions\n", len(completed)-count)
					break
				}

				fmt.Printf("%s - %s\n", p.Decision.Timestamp.Format("2006-01-02"), p.Title)
				fmt.Printf("  ID: %s\n", p.ID)
				fmt.Printf("  Decision: %s\n", p.Decision.Result)
				fmt.Printf("  Rationale: %s\n", p.Decision.Rationale)
				fmt.Printf("  Consultations: %d\n", len(p.Consultations))
				fmt.Println()

				count++
			}

			// Show decision statistics
			decisionStats := make(map[proposal.DecisionResult]int)
			for _, p := range completed {
				if p.Decision != nil {
					decisionStats[p.Decision.Result]++
				}
			}

			fmt.Println("Decision Statistics:")
			fmt.Printf("  Approved: %d\n", decisionStats[proposal.DecisionApproved])
			fmt.Printf("  Rejected: %d\n", decisionStats[proposal.DecisionRejected])
			fmt.Printf("  Deferred: %d\n", decisionStats[proposal.DecisionDeferred])
			fmt.Printf("  No Consensus: %d\n", decisionStats[proposal.DecisionNoConsensus])

			return nil
		},
	}

	cmd.Flags().IntVarP(&days, "days", "d", 30, "Number of days of history to show")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of decisions to show")

	return cmd
}

// newStatusHealthCmd shows collective health metrics
func newStatusHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Show collective health metrics",
		Long: `Display metrics about the health of the collective's
decision-making processes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all proposals for analysis
			allProposals, err := proposal.List(proposal.ListFilter{ShowAll: true})
			if err != nil {
				return fmt.Errorf("failed to get proposals: %w", err)
			}

			fmt.Println("Collective Health Metrics")
			fmt.Println("========================")

			// Calculate metrics
			totalProposals := len(allProposals)
			if totalProposals == 0 {
				fmt.Println("No proposals yet - the collective is just beginning!")
				return nil
			}

			// Count decisions
			approved := 0
			rejected := 0
			noConsensus := 0
			avgConsultations := 0
			totalConsultations := 0

			for _, p := range allProposals {
				if p.Decision != nil {
					switch p.Decision.Result {
					case proposal.DecisionApproved:
						approved++
					case proposal.DecisionRejected:
						rejected++
					case proposal.DecisionNoConsensus:
						noConsensus++
					}
				}
				totalConsultations += len(p.Consultations)
			}

			if totalProposals > 0 {
				avgConsultations = totalConsultations / totalProposals
			}

			// Display metrics
			fmt.Printf("\nTotal Proposals: %d\n", totalProposals)
			fmt.Printf("Consensus Success Rate: %.1f%%\n", float64(approved)*100/float64(totalProposals))
			fmt.Printf("Average Consultations per Proposal: %d\n", avgConsultations)

			// Participation health
			fmt.Println("\nParticipation Health:")
			if avgConsultations < 2 {
				fmt.Println("  ⚠️  Low participation - consider encouraging more input")
			} else if avgConsultations < 5 {
				fmt.Println("  ✓ Moderate participation")
			} else {
				fmt.Println("  ✓ High participation - healthy collective engagement!")
			}

			// Consensus health
			fmt.Println("\nConsensus Health:")
			consensusRate := float64(approved) / float64(approved+rejected+noConsensus)
			if consensusRate < 0.5 {
				fmt.Println("  ⚠️  Low consensus rate - consider more dialogue")
			} else if consensusRate < 0.8 {
				fmt.Println("  ✓ Moderate consensus rate")
			} else {
				fmt.Println("  ✓ High consensus rate - strong collective alignment!")
			}

			// Blocking concerns
			blocked := 0
			for _, p := range allProposals {
				if p.Status == proposal.StatusBlocked {
					blocked++
				}
			}

			fmt.Println("\nBlocked Proposals:")
			if blocked == 0 {
				fmt.Println("  ✓ No blocked proposals - concerns are being addressed")
			} else {
				fmt.Printf("  ⚠️  %d blocked proposals need attention\n", blocked)
			}

			// Recent activity
			recentProposals := 0
			for _, p := range allProposals {
				if time.Since(p.Date) < 30*24*time.Hour {
					recentProposals++
				}
			}

			fmt.Println("\nActivity Level:")
			if recentProposals == 0 {
				fmt.Println("  ⚠️  No recent activity - collective may be dormant")
			} else if recentProposals < 5 {
				fmt.Println("  ✓ Low but steady activity")
			} else {
				fmt.Println("  ✓ High activity - engaged collective!")
			}

			return nil
		},
	}

	return cmd
}

// showProposalSummary displays a brief summary of a proposal
func showProposalSummary(p *proposal.Proposal) {
	urgencyIcon := ""
	switch p.Urgency {
	case proposal.UrgencyEmergency:
		urgencyIcon = "🚨"
	case proposal.UrgencyHigh:
		urgencyIcon = "⚠️ "
	case proposal.UrgencyMedium:
		urgencyIcon = "📌"
	case proposal.UrgencyLow:
		urgencyIcon = "💭"
	}

	fmt.Printf("  %s %s - %s\n", urgencyIcon, p.ID, p.Title)
	fmt.Printf("     Urgency: %s | Proposer: %s | Date: %s\n",
		p.Urgency, p.Proposer, p.Date.Format("2006-01-02"))
	
	if len(p.Consultations) > 0 {
		fmt.Printf("     Consultations: %d", len(p.Consultations))
		concerns := p.GetBlockingConcerns()
		if len(concerns) > 0 {
			fmt.Printf(" (⚠️  %d concerns)", len(concerns))
		}
		fmt.Println()
	}
	
	if p.ConsensusStatus != "" {
		fmt.Printf("     Status: %s\n", p.ConsensusStatus)
	}
}