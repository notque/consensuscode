// Package web provides the embedded web interface for CollectiveFlow.
// Following horizontal principles, the web interface has no admin panels,
// no special user roles, and presents information in a transparent,
// community-oriented way.
package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"sort"
	"time"

	"collectiveflow/internal/proposal"
	"collectiveflow/internal/storage"
)

//go:embed templates/* static/*
var embeddedFS embed.FS

// Server represents the web server for CollectiveFlow
type Server struct {
	store     storage.ProposalStore
	templates *template.Template
	addr      string
}

// Config holds configuration for the web server
type Config struct {
	Store storage.ProposalStore
	Addr  string // e.g., ":8080"
}

// NewServer creates a new web server instance
func NewServer(cfg Config) (*Server, error) {
	s := &Server{
		store: cfg.Store,
		addr:  cfg.Addr,
	}

	// Load templates with custom functions
	tmpl, err := s.loadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}
	s.templates = tmpl

	return s, nil
}

// loadTemplates loads all HTML templates with custom functions
func (s *Server) loadTemplates() (*template.Template, error) {
	// Custom template functions
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"timeAgo": func(t time.Time) string {
			duration := time.Since(t)
			switch {
			case duration < time.Minute:
				return "just now"
			case duration < time.Hour:
				return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
			case duration < 24*time.Hour:
				return fmt.Sprintf("%d hours ago", int(duration.Hours()))
			case duration < 7*24*time.Hour:
				return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
			default:
				return t.Format("2006-01-02")
			}
		},
		"statusClass": func(status proposal.ProposalStatus) string {
			switch status {
			case proposal.StatusProposed:
				return "status-proposed"
			case proposal.StatusConsultation:
				return "status-consultation"
			case proposal.StatusConsensus:
				return "status-consensus"
			case proposal.StatusImplemented:
				return "status-implemented"
			case proposal.StatusWithdrawn:
				return "status-withdrawn"
			case proposal.StatusBlocked:
				return "status-blocked"
			default:
				return "status-unknown"
			}
		},
		"urgencyClass": func(urgency proposal.UrgencyLevel) string {
			switch urgency {
			case proposal.UrgencyEmergency:
				return "urgency-emergency"
			case proposal.UrgencyHigh:
				return "urgency-high"
			case proposal.UrgencyMedium:
				return "urgency-medium"
			case proposal.UrgencyLow:
				return "urgency-low"
			default:
				return "urgency-unknown"
			}
		},
		"percentage": func(count, total int) int {
			if total == 0 {
				return 0
			}
			return (count * 100) / total
		},
	}

	// Load templates from embedded filesystem
	tmplFS, err := fs.Sub(embeddedFS, "templates")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(tmplFS, "*.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// Start starts the web server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve static files
	staticFS, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		return fmt.Errorf("failed to load static files: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Register routes
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/proposals", s.handleProposals)
	mux.HandleFunc("/proposals/", s.handleProposal)
	mux.HandleFunc("/collective", s.handleCollective)
	mux.HandleFunc("/about", s.handleAbout)

	fmt.Printf("\n🌐 CollectiveFlow Web Interface\n")
	fmt.Printf("================================\n")
	fmt.Printf("Server starting at http://localhost%s\n", s.addr)
	fmt.Printf("\nHorizontal principles in action:\n")
	fmt.Printf("  • No admin panels or privileged access\n")
	fmt.Printf("  • All proposals visible to everyone\n")
	fmt.Printf("  • Transparent consensus process\n")
	fmt.Printf("  • Community bulletin board design\n")
	fmt.Printf("\nPress Ctrl+C to stop the server\n\n")

	return http.ListenAndServe(s.addr, mux)
}

// handleIndex handles the homepage
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Load all proposals for overview
	proposalsInterface, err := s.store.ListAll()
	if err != nil {
		http.Error(w, "Failed to load proposals", http.StatusInternalServerError)
		return
	}

	proposals := make([]*proposal.Proposal, 0, len(proposalsInterface))
	for _, p := range proposalsInterface {
		if prop, ok := p.(*proposal.Proposal); ok {
			proposals = append(proposals, prop)
		}
	}

	// Calculate collective statistics (non-hierarchical metrics)
	stats := calculateStats(proposals)

	data := map[string]interface{}{
		"Title":      "CollectiveFlow",
		"ActiveTab":  "home",
		"Stats":      stats,
		"Recent":     getRecentProposals(proposals, 5),
		"NeedsInput": getNeedsInputProposals(proposals),
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleProposals handles the proposals list page
func (s *Server) handleProposals(w http.ResponseWriter, r *http.Request) {
	// Get filter parameters
	statusFilter := r.URL.Query().Get("status")
	urgencyFilter := r.URL.Query().Get("urgency")

	// Load all proposals
	proposalsInterface, err := s.store.ListAll()
	if err != nil {
		http.Error(w, "Failed to load proposals", http.StatusInternalServerError)
		return
	}

	proposals := make([]*proposal.Proposal, 0, len(proposalsInterface))
	for _, p := range proposalsInterface {
		if prop, ok := p.(*proposal.Proposal); ok {
			// Apply filters
			if statusFilter != "" && string(prop.Status) != statusFilter {
				continue
			}
			if urgencyFilter != "" && string(prop.Urgency) != urgencyFilter {
				continue
			}
			proposals = append(proposals, prop)
		}
	}

	// Sort by date (newest first)
	sort.Slice(proposals, func(i, j int) bool {
		return proposals[i].Date.After(proposals[j].Date)
	})

	data := map[string]interface{}{
		"Title":         "Proposals - CollectiveFlow",
		"ActiveTab":     "proposals",
		"Proposals":     proposals,
		"StatusFilter":  statusFilter,
		"UrgencyFilter": urgencyFilter,
	}

	if err := s.templates.ExecuteTemplate(w, "proposals.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleProposal handles individual proposal pages
func (s *Server) handleProposal(w http.ResponseWriter, r *http.Request) {
	// Extract proposal ID from URL path
	id := r.URL.Path[len("/proposals/"):]
	if id == "" {
		http.Redirect(w, r, "/proposals", http.StatusFound)
		return
	}

	// Load proposal
	propInterface, err := s.store.Load(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	prop, ok := propInterface.(*proposal.Proposal)
	if !ok {
		http.Error(w, "Invalid proposal data", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":     prop.Title + " - CollectiveFlow",
		"ActiveTab": "proposals",
		"Proposal":  prop,
	}

	if err := s.templates.ExecuteTemplate(w, "proposal.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCollective handles the collective status page
func (s *Server) handleCollective(w http.ResponseWriter, r *http.Request) {
	// Load all proposals
	proposalsInterface, err := s.store.ListAll()
	if err != nil {
		http.Error(w, "Failed to load proposals", http.StatusInternalServerError)
		return
	}

	proposals := make([]*proposal.Proposal, 0, len(proposalsInterface))
	for _, p := range proposalsInterface {
		if prop, ok := p.(*proposal.Proposal); ok {
			proposals = append(proposals, prop)
		}
	}

	stats := calculateStats(proposals)

	data := map[string]interface{}{
		"Title":     "Collective Status - CollectiveFlow",
		"ActiveTab": "collective",
		"Stats":     stats,
	}

	if err := s.templates.ExecuteTemplate(w, "collective.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAbout handles the about page
func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":     "About - CollectiveFlow",
		"ActiveTab": "about",
	}

	if err := s.templates.ExecuteTemplate(w, "about.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Stats represents collective statistics
type Stats struct {
	Total         int
	Active        int
	Implemented   int
	ConsensusRate int
	ByStatus      map[string]int
	ByUrgency     map[string]int
}

// calculateStats computes non-hierarchical collective metrics
func calculateStats(proposals []*proposal.Proposal) Stats {
	stats := Stats{
		Total:     len(proposals),
		ByStatus:  make(map[string]int),
		ByUrgency: make(map[string]int),
	}

	consensus := 0
	for _, p := range proposals {
		// Count by status
		stats.ByStatus[string(p.Status)]++

		// Count active proposals (not terminal states)
		if p.Status != proposal.StatusImplemented &&
		   p.Status != proposal.StatusWithdrawn &&
		   p.Status != proposal.StatusBlocked {
			stats.Active++
		}

		// Count implemented
		if p.Status == proposal.StatusImplemented {
			stats.Implemented++
		}

		// Count reached consensus (for rate calculation)
		if p.Status == proposal.StatusConsensus || p.Status == proposal.StatusImplemented {
			consensus++
		}

		// Count by urgency
		stats.ByUrgency[string(p.Urgency)]++
	}

	// Calculate consensus rate (percentage)
	if stats.Total > 0 {
		stats.ConsensusRate = (consensus * 100) / stats.Total
	}

	return stats
}

// getRecentProposals returns the N most recent proposals
func getRecentProposals(proposals []*proposal.Proposal, n int) []*proposal.Proposal {
	// Sort by date (newest first)
	sorted := make([]*proposal.Proposal, len(proposals))
	copy(sorted, proposals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})

	if len(sorted) > n {
		return sorted[:n]
	}
	return sorted
}

// getNeedsInputProposals returns proposals in consultation that need input
func getNeedsInputProposals(proposals []*proposal.Proposal) []*proposal.Proposal {
	var needsInput []*proposal.Proposal
	for _, p := range proposals {
		if p.Status == proposal.StatusConsultation {
			needsInput = append(needsInput, p)
		}
	}

	// Sort by urgency (emergency first)
	sort.Slice(needsInput, func(i, j int) bool {
		urgencyOrder := map[proposal.UrgencyLevel]int{
			proposal.UrgencyEmergency: 0,
			proposal.UrgencyHigh:      1,
			proposal.UrgencyMedium:    2,
			proposal.UrgencyLow:       3,
		}
		return urgencyOrder[needsInput[i].Urgency] < urgencyOrder[needsInput[j].Urgency]
	})

	return needsInput
}
