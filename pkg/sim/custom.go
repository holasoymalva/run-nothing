package sim

import (
	"fmt"
	"math/rand"
	"strings"
)

// BuildCustomScenario synthesizes a realistic scenario tailored to a custom user prompt
func BuildCustomScenario(userPrompt string) Scenario {
	// Pick a base scenario and customize it to the user's prompt
	base := Scenarios[rand.Intn(len(Scenarios))]

	slug := slugify(userPrompt)
	if slug == "" {
		slug = "custom-task"
	}

	targetFile := fmt.Sprintf("internal/%s/handler.go", slug)

	return Scenario{
		Prompt: userPrompt,
		ScoutActions: []struct{ Tool, Args, Result string }{
			{"Grep", fmt.Sprintf(`"%s" path="."`, extractKeyword(userPrompt)), "Found 11 matches across 4 files"},
			{"Read", fmt.Sprintf(`"%s" (lines 1-120)`, targetFile), "Read 120 lines (4.1 KB)"},
		},
		Thoughts: []string{
			fmt.Sprintf("Analyzing requirements for: %s...", truncate(userPrompt, 42)),
			"Inspecting call graph and interface boundaries...",
			"Synthesizing zero-allocation implementation with strict error handling...",
		},
		Diffs: []DiffBlock{
			{
				Filename: targetFile,
				Lines: []string{
					"@@ -28,7 +28,12 @@ func ExecuteTask(ctx context.Context, req *Request) (*Response, error) {",
					"-\tif req == nil {",
					"-\t\treturn nil, ErrEmptyRequest",
					"+\tif err := req.Validate(ctx); err != nil {",
					"+\t\treturn nil, fmt.Errorf(\"validation failed: %w\", err)",
					" \t}",
					"+\tmetrics.RecordInvocation(req.OperationName)",
					" \treturn processOptimized(ctx, req)",
				},
			},
		},
		TestCommand: fmt.Sprintf("go test -v -race ./internal/%s/...", slug),
		TestOutputs: []string{
			fmt.Sprintf("=== RUN   TestExecuteTask_%s", toCamel(slug)),
			fmt.Sprintf("--- PASS: TestExecuteTask_%s (0.04s)", toCamel(slug)),
			"PASS",
			fmt.Sprintf("ok  \tgithub.com/acme/core/internal/%s\t0.184s (coverage: 95.8%% of statements)", slug),
		},
		ReviewNotes: base.ReviewNotes,
		CommitMessage: fmt.Sprintf("feat(%s): implement %s", slug, truncate(strings.ToLower(userPrompt), 48)),
		BranchName:    fmt.Sprintf("claude/%s", slug),
		DeployTarget:  base.DeployTarget,
		Conclusion:    fmt.Sprintf("Completed the requested changes for **%s** in `%s`. All unit tests and race detector checks pass with 95.8%% statement coverage.", userPrompt, targetFile),
	}
}

func extractKeyword(prompt string) string {
	words := strings.Fields(prompt)
	for _, w := range words {
		clean := strings.Trim(w, ".,!?\"'`")
		if len(clean) > 4 {
			return clean
		}
	}
	return "Handler"
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	words := strings.Fields(s)
	count := 0
	for _, w := range words {
		clean := ""
		for _, r := range w {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				clean += string(r)
			}
		}
		if len(clean) > 2 {
			if count > 0 {
				b.WriteRune('-')
			}
			b.WriteString(clean)
			count++
			if count >= 2 {
				break
			}
		}
	}
	return b.String()
}

func toCamel(s string) string {
	parts := strings.Split(s, "-")
	var b strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
