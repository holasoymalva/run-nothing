package sim

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/holasoymalva/run-nothing/pkg/stages"
	"github.com/holasoymalva/run-nothing/pkg/ui"
)

// EngineOptions configures the simulation run
type EngineOptions struct {
	Speed           string // "normal", "fast", "cinematic", "instant"
	Model           string
	Filter          *stages.FilterConfig
	InteractiveMode bool
}

// GetDelay returns the scaled duration according to current speed
func (opts *EngineOptions) GetDelay(base time.Duration) time.Duration {
	switch opts.Speed {
	case "instant":
		return 0
	case "fast":
		return base / 3
	case "cinematic":
		return base * 2
	case "normal":
		return base
	default:
		return base
	}
}

// Sleep pauses for the scaled duration
func (opts *EngineOptions) Sleep(base time.Duration) {
	d := opts.GetDelay(base)
	if d > 0 {
		time.Sleep(d)
	}
}

// RunScenario executes a scenario through all enabled stages
func RunScenario(ctx context.Context, sc Scenario, opts EngineOptions) error {
	startTime := time.Now()

	// 1. Prompt header
	fmt.Printf("%s╭─ %sUser%s %s───────────────────────────────────────────────╮%s\n",
		ui.DarkGray, ui.Bold, ui.Reset, ui.DarkGray, ui.Reset)
	fmt.Printf("%s│%s %s%s%s\n", ui.DarkGray, ui.Reset, ui.White, sc.Prompt, ui.Reset)
	fmt.Printf("%s╰──────────────────────────────────────────────────────────╯%s\n\n",
		ui.DarkGray, ui.Reset)

	opts.Sleep(400 * time.Millisecond)

	// Stage: Scout (Codebase exploration)
	if opts.Filter.IsEnabled("scout") {
		for _, action := range sc.ScoutActions {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			ui.PrintToolCall(action.Tool, action.Args)
			opts.Sleep(300 * time.Millisecond)
			fmt.Printf("  %s└─ %s%s\n", ui.Dim, action.Result, ui.Reset)
			opts.Sleep(200 * time.Millisecond)
		}
		fmt.Println()
	}

	// Stage: Think (Chain of thought reasoning)
	if opts.Filter.IsEnabled("think") {
		spinner := ui.StartSpinner("Thinking...")
		for _, thought := range sc.Thoughts {
			select {
			case <-ctx.Done():
				spinner.Stop("", "")
				return ctx.Err()
			default:
			}
			opts.Sleep(600 * time.Millisecond)
			spinner.UpdateLabel(fmt.Sprintf("Thinking... %s(%s)%s", ui.Gray, thought, ui.Reset))
			opts.Sleep(700 * time.Millisecond)
		}
		spinner.Stop(ui.Green+"✔"+ui.Reset, ui.Bold+"Formulated execution plan"+ui.Reset)
		fmt.Println()
	}

	// Stage: Codegen (Unified file diffs)
	if opts.Filter.IsEnabled("codegen") {
		for _, diff := range sc.Diffs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			ui.PrintToolCall("EditFile", fmt.Sprintf(`path="%s"`, diff.Filename))
			ui.PrintUnifiedDiff(diff.Filename, diff.Lines)
			opts.Sleep(400 * time.Millisecond)
		}
	}

	// Stage: Test (Simulated test runners)
	if opts.Filter.IsEnabled("test") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		ui.PrintToolCall("Bash", sc.TestCommand)
		opts.Sleep(300 * time.Millisecond)
		for _, outLine := range sc.TestOutputs {
			fmt.Printf("  %s%s%s\n", ui.Dim, outLine, ui.Reset)
			opts.Sleep(60 * time.Millisecond)
		}
		fmt.Printf("  %s✔ All checks passed!%s\n\n", ui.Green+ui.Bold, ui.Reset)
	}

	// Stage: Review (Multi-agent review)
	if opts.Filter.IsEnabled("review") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		fmt.Printf("%s● UltraReview (Multi-Agent Verification):%s\n", ui.Purple, ui.Reset)
		for _, note := range sc.ReviewNotes {
			fmt.Printf("  %s✓%s %s\n", ui.Green, ui.Reset, note)
			opts.Sleep(150 * time.Millisecond)
		}
		fmt.Println()
	}

	// Stage: Git (Atomic commit & branch)
	if opts.Filter.IsEnabled("git") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		ui.PrintToolCall("Bash", fmt.Sprintf(`git commit -m "%s"`, sc.CommitMessage))
		fmt.Printf("  %s[detached HEAD %07x] %s%s\n", ui.Dim, rand.Int31n(0x0fffffff), sc.CommitMessage, ui.Reset)
		fmt.Printf("  %sBranch '%s' updated.%s\n\n", ui.Cyan, sc.BranchName, ui.Reset)
		opts.Sleep(250 * time.Millisecond)
	}

	// Stage: Deploy (Simulated rollout progress)
	if opts.Filter.IsEnabled("deploy") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		fmt.Printf("%s● Canary Deployment: %s%s\n", ui.SoftCyan, sc.DeployTarget, ui.Reset)
		totalSteps := 10
		for step := 1; step <= totalSteps; step++ {
			ui.RenderProgressBar(step, totalSteps, 28, "Syncing pods")
			opts.Sleep(90 * time.Millisecond)
		}
		fmt.Printf("  %s✔ Rollout complete! Zero 5xx errors recorded.%s\n\n", ui.Green, ui.Reset)
	}

	// Claude Conclusion
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fmt.Printf("%sClaude:%s\n", ui.Coral+ui.Bold, ui.Reset)
	if opts.Speed == "instant" {
		fmt.Println(sc.Conclusion)
	} else {
		charDelay := opts.GetDelay(12 * time.Millisecond)
		ui.Typewriter(sc.Conclusion, charDelay)
	}
	fmt.Println()

	// Cost and token turn summary
	duration := time.Since(startTime)
	simulatedCost := 0.025 + rand.Float64()*0.035
	inTokens := 12000 + rand.Intn(18000)
	outTokens := 850 + rand.Intn(1200)

	ui.PrintTurnSummary(simulatedCost, inTokens, outTokens, duration)
	return nil
}
