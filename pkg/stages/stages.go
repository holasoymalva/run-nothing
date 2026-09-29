package stages

import (
	"fmt"
	"strings"
)

// Stage represents a simulation step in Claude Code's agentic workflow
type Stage struct {
	ID          string
	Name        string
	Description string
}

// AvailableStages lists all stages that run-nothing can simulate
var AvailableStages = []Stage{
	{
		ID:          "scout",
		Name:        "Codebase Exploration",
		Description: "Indexing repository, grepping symbols, AST parsing, and reading context",
	},
	{
		ID:          "think",
		Name:        "Deep Reasoning",
		Description: "Extended thinking spinner and architectural hypothesis formulation",
	},
	{
		ID:          "codegen",
		Name:        "Code Synthesis & Diffs",
		Description: "Generating refactors, bugfixes, and colorized unified git diffs",
	},
	{
		ID:          "test",
		Name:        "Test Runner & Verification",
		Description: "Simulating test execution, linters, race detectors, and benchmarks",
	},
	{
		ID:          "review",
		Name:        "Multi-Agent Code Review",
		Description: "Automated UltraReview checking for security flaws, regressions, and style",
	},
	{
		ID:          "git",
		Name:        "Git Version Control",
		Description: "Simulating atomic staging, conventional commit messages, and PR creation",
	},
	{
		ID:          "deploy",
		Name:        "CI/CD & Canary Rollout",
		Description: "Simulating container packaging, pipeline execution, and cluster health",
	},
}

// FilterConfig holds stage inclusion and exclusion criteria
type FilterConfig struct {
	Included map[string]bool
	Excluded map[string]bool
}

// NewFilterConfig parses included and excluded stage slices
func NewFilterConfig(includes, excludes []string) *FilterConfig {
	cfg := &FilterConfig{
		Included: make(map[string]bool),
		Excluded: make(map[string]bool),
	}

	for _, s := range includes {
		s = strings.TrimSpace(strings.ToLower(s))
		if s != "" {
			// handle comma separated strings
			parts := strings.Split(s, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					cfg.Included[p] = true
				}
			}
		}
	}

	for _, s := range excludes {
		s = strings.TrimSpace(strings.ToLower(s))
		if s != "" {
			parts := strings.Split(s, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					cfg.Excluded[p] = true
				}
			}
		}
	}

	return cfg
}

// IsEnabled checks whether a given stage ID should run
func (cfg *FilterConfig) IsEnabled(stageID string) bool {
	stageID = strings.ToLower(stageID)

	// If explicitly excluded, do not run
	if cfg.Excluded[stageID] {
		return false
	}

	// If specific stages were included, only run those
	if len(cfg.Included) > 0 {
		return cfg.Included[stageID]
	}

	// By default, run all non-excluded stages
	return true
}

// PrintAvailableStages prints the list of stages in a clean format
func PrintAvailableStages() {
	fmt.Println("\nAvailable simulation stages:")
	for _, s := range AvailableStages {
		fmt.Printf("  • %-10s %s (%s)\n", s.ID, s.Name, s.Description)
	}
	fmt.Println()
}
