package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/holasoymalva/run-nothing/pkg/sim"
	"github.com/holasoymalva/run-nothing/pkg/stages"
	"github.com/holasoymalva/run-nothing/pkg/ui"
)

// RunInteractiveLoop starts the interactive Claude Code REPL
func RunInteractiveLoop(ctx context.Context, opts sim.EngineOptions, version string) {
	scanner := bufio.NewScanner(os.Stdin)
	var totalSimulatedCost float64 = 0.0

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "run-nothing"
	}

	fmt.Printf("%scwd: %s%s %s(git: main)%s\n\n", ui.Dim, cwd, ui.Reset, ui.Cyan, ui.Reset)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n%sSession terminated.%s\n", ui.Coral, ui.Reset)
			return
		default:
		}

		// Prompt box
		fmt.Printf("%s> %s", ui.Coral+ui.Bold, ui.Reset)

		if !scanner.Scan() {
			fmt.Println()
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Handle slash commands
		if strings.HasPrefix(line, "/") {
			handled := handleSlashCommand(line, &opts, &totalSimulatedCost, version)
			if handled == "exit" {
				break
			}
			continue
		}

		// Execute simulated scenario based on user prompt
		scenario := sim.BuildCustomScenario(line)
		err := sim.RunScenario(ctx, scenario, opts)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			fmt.Printf("%sError: %v%s\n", ui.Red, err, ui.Reset)
		}

		totalSimulatedCost += 0.038
	}

	fmt.Printf("\n%sGoodbye!%s\n", ui.Coral+ui.Bold, ui.Reset)
	fmt.Printf("%sTotal actual cost: $0.00 (run-nothing runs completely locally & free)%s\n", ui.Dim, ui.Reset)
}

func handleSlashCommand(cmd string, opts *sim.EngineOptions, totalCost *float64, version string) string {
	parts := strings.Fields(cmd)
	switch strings.ToLower(parts[0]) {
	case "/exit", "/quit", "/q":
		return "exit"

	case "/help", "/h":
		fmt.Println("\nClaude Code Interactive Commands:")
		fmt.Printf("  %s/help%s         Show this help message\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/cost%s         Display session token usage and estimated cost\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/model [name]%s Switch active simulated LLM\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/compact%s      Simulate context window compaction\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/doctor%s       Run simulated system health checks\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/stages%s       List all available simulation stages\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/clear%s        Clear screen\n", ui.Cyan, ui.Reset)
		fmt.Printf("  %s/exit%s         Exit Claude Code\n\n", ui.Cyan, ui.Reset)

	case "/cost":
		fmt.Printf("\n%s● Session Cost Tracker:%s\n", ui.Coral, ui.Reset)
		fmt.Printf("  Simulated API cost:  %s$%0.4f%s\n", ui.Yellow, *totalCost, ui.Reset)
		fmt.Printf("  Actual dollar cost:  %s$0.00%s (run-nothing runs nothing)\n\n", ui.Green+ui.Bold, ui.Reset)

	case "/model":
		if len(parts) > 1 {
			opts.Model = parts[1]
			fmt.Printf("%sActive model switched to %s%s%s\n\n", ui.Green, ui.Bold, opts.Model, ui.Reset)
		} else {
			fmt.Printf("Current model: %s%s%s (Options: claude-3-7-sonnet, claude-3-5-sonnet, claude-4-opus, claude-3-5-haiku)\n\n",
				ui.Bold, opts.Model, ui.Reset)
		}

	case "/compact":
		fmt.Printf("%sCompacting conversation context window...%s\n", ui.Dim, ui.Reset)
		for i := 1; i <= 8; i++ {
			ui.RenderProgressBar(i, 8, 24, "Pruning history")
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Printf("%s✔ Context reduced by 64.2%% (128k -> 46k tokens)%s\n\n", ui.Green, ui.Reset)

	case "/doctor":
		osName := runtime.GOOS
		if term := os.Getenv("TERM_PROGRAM"); term != "" {
			osName += fmt.Sprintf(" (%s)", term)
		}
		fmt.Printf("\n%sClaude Code Health Check (Simulated):%s\n", ui.Coral+ui.Bold, ui.Reset)
		fmt.Printf("  %s✓%s OS: %s (%s/%s)\n", ui.Green, ui.Reset, osName, runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  %s✓%s Claude Code Native: v%s\n", ui.Green, ui.Reset, version)
		fmt.Printf("  %s✓%s Git: Clean working tree detected\n", ui.Green, ui.Reset)
		fmt.Printf("  %s✓%s Network Sandbox: 100%% airgapped safe\n", ui.Green, ui.Reset)
		fmt.Printf("  %s✓%s API Token: Simulated unlimited enterprise quota\n\n", ui.Green, ui.Reset)

	case "/stages":
		stages.PrintAvailableStages()

	case "/clear":
		ui.ClearScreen()
		ui.PrintBanner(version, opts.Model)

	default:
		fmt.Printf("%sUnknown command '%s'. Type /help for available commands.%s\n\n", ui.Red, cmd, ui.Reset)
	}
	return ""
}
