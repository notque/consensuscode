package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"collectiveflow/internal/proposal"
)

// knownAgents lists all agents in the collective. This is derived from the
// agents/ directory in the repository. If the collective grows, update this
// list through a consensus proposal.
var knownAgents = []string{
	"api-design-specialist",
	"consensus-coordinator",
	"database-design-specialist",
	"david-graeber-agent",
	"devops-local-infrastructure",
	"documentation-specialist",
	"flask-web-developer",
	"frontend-specialist",
	"go-code-quality-specialist",
	"go-systems-developer",
	"noam-chomsky-agent",
	"product-steward",
	"python-testing-specialist",
	"ux-research-specialist",
	"web-security-specialist",
}

// newDashboardCmd creates the dashboard command
func newDashboardCmd() *cobra.Command {
	var days int

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Show a collective activity summary",
		Long: `Display an at-a-glance summary of the collective's proposal activity.

The dashboard shows:
  - Total proposals grouped by status
  - Recent activity within a configurable window
  - Agents who have not yet provided input on active proposals

This helps every member of the collective stay informed without
creating hierarchical oversight.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDashboard(days)
		},
	}

	cmd.Flags().IntVarP(&days, "days", "d", 30, "Number of days to include in the recent-activity window")

	return cmd
}

// runDashboard produces the dashboard output.
func runDashboard(days int) error {
	allProposals, err := proposal.List(proposal.ListFilter{ShowAll: true})
	if err != nil {
		return fmt.Errorf("could not load proposals: %w\n\nTroubleshooting:\n  - Make sure the data/proposals/ directory exists\n  - Check file permissions on the proposal YAML files", err)
	}

	printDashboardHeader()
	printStatusSummary(allProposals)
	printRecentActivity(allProposals, days)
	printMissingInput(allProposals)
	printDashboardFooter()

	return nil
}

// printDashboardHeader prints the dashboard banner.
func printDashboardHeader() {
	fmt.Println()
	fmt.Println("===================================")
	fmt.Println("  CollectiveFlow Dashboard")
	fmt.Printf("  %s\n", time.Now().Format("2006-01-02 15:04"))
	fmt.Println("===================================")
}

// printStatusSummary shows total proposals grouped by status.
func printStatusSummary(proposals []*proposal.Proposal) {
	fmt.Println()
	fmt.Printf("Proposals: %d total\n", len(proposals))

	if len(proposals) == 0 {
		fmt.Println("  (none yet)")
		return
	}

	counts := make(map[proposal.ProposalStatus]int)
	for _, p := range proposals {
		counts[p.Status]++
	}

	// Print in a meaningful order: active states first, then terminal.
	order := []struct {
		status proposal.ProposalStatus
		label  string
	}{
		{proposal.StatusConsultation, "In consultation"},
		{proposal.StatusProposed, "Proposed"},
		{proposal.StatusConsensus, "Consensus reached"},
		{proposal.StatusBlocked, "Blocked"},
		{proposal.StatusImplemented, "Implemented"},
		{proposal.StatusWithdrawn, "Withdrawn"},
	}

	for _, o := range order {
		n := counts[o.status]
		if n > 0 {
			fmt.Printf("  %-20s %d\n", o.label, n)
		}
	}
}

// printRecentActivity shows proposals with activity in the last N days.
func printRecentActivity(proposals []*proposal.Proposal, days int) {
	fmt.Println()
	fmt.Printf("Recent activity (last %d days):\n", days)

	cutoff := time.Now().AddDate(0, 0, -days)

	type activityEntry struct {
		date    time.Time
		id      string
		title   string
		event   string
	}

	var entries []activityEntry

	for _, p := range proposals {
		// Check proposal creation date.
		if p.Date.After(cutoff) {
			entries = append(entries, activityEntry{
				date:  p.Date,
				id:    p.ID,
				title: p.Title,
				event: "created",
			})
		}

		// Check consensus history for recent events.
		for _, h := range p.ConsensusHistory {
			if h.Timestamp.After(cutoff) && h.Event != "proposal_created" {
				label := h.Event
				switch h.Event {
				case "status_changed":
					label = "status changed"
				case "consultation_received":
					label = fmt.Sprintf("input from %s", h.Actor)
				case "decision_recorded":
					label = "decision recorded"
				case "proposal_updated":
					label = "updated"
				}
				entries = append(entries, activityEntry{
					date:  h.Timestamp,
					id:    p.ID,
					title: p.Title,
					event: label,
				})
			}
		}
	}

	if len(entries) == 0 {
		fmt.Printf("  No activity in the last %d days.\n", days)
		return
	}

	// Sort newest first.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].date.After(entries[j].date)
	})

	// Cap at 15 to keep the dashboard readable.
	shown := entries
	if len(shown) > 15 {
		shown = shown[:15]
	}

	for _, e := range shown {
		truncTitle := e.title
		if len(truncTitle) > 50 {
			truncTitle = truncTitle[:47] + "..."
		}
		fmt.Printf("  %s  %-18s  %s\n",
			e.date.Format("Jan 02 15:04"),
			e.event,
			truncTitle,
		)
	}

	if len(entries) > 15 {
		fmt.Printf("  ... and %d more events\n", len(entries)-15)
	}
}

// printMissingInput identifies agents who have not provided input on active proposals.
func printMissingInput(proposals []*proposal.Proposal) {
	fmt.Println()
	fmt.Println("Agents missing input on active proposals:")

	// Collect active proposals (in consultation).
	var active []*proposal.Proposal
	for _, p := range proposals {
		if p.Status == proposal.StatusConsultation {
			active = append(active, p)
		}
	}

	if len(active) == 0 {
		fmt.Println("  No proposals currently in consultation.")
		return
	}

	// For each active proposal, find which known agents have not contributed.
	anyMissing := false
	for _, p := range active {
		contributed := make(map[string]bool)
		for _, c := range p.Consultations {
			contributed[c.Contributor] = true
		}

		var missing []string
		for _, agent := range knownAgents {
			if !contributed[agent] {
				missing = append(missing, agent)
			}
		}

		if len(missing) > 0 {
			anyMissing = true
			truncTitle := p.Title
			if len(truncTitle) > 50 {
				truncTitle = truncTitle[:47] + "..."
			}
			fmt.Printf("\n  %s (%s):\n", truncTitle, p.ID)
			fmt.Printf("    Missing: %s\n", strings.Join(missing, ", "))
		}
	}

	if !anyMissing {
		fmt.Println("  All known agents have provided input on active proposals.")
	}
}

// printDashboardFooter prints the closing hint.
func printDashboardFooter() {
	fmt.Println()
	fmt.Println("---")
	fmt.Println("Run 'collectiveflow status active' for full details on active proposals.")
	fmt.Println("Run 'collectiveflow proposal list --all' to see every proposal.")
	fmt.Println()
}
