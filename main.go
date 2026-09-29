package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/holasoymalva/run-nothing/pkg/repl"
	"github.com/holasoymalva/run-nothing/pkg/sim"
	"github.com/holasoymalva/run-nothing/pkg/stages"
	"github.com/holasoymalva/run-nothing/pkg/ui"
)

const Version = "2.1.139"

type cliConfig struct {
	showHelp        bool
	showVersion     bool
	interactive     bool
	once            bool
	skipPermissions bool
	speed           string
	model           string
	prompt          string
	includedStages  []string
	excludedStages  []string
}

func parseArgs(args []string) cliConfig {
	cfg := cliConfig{
		speed: "normal",
		model: "claude-sonnet-4-6",
	}

	knownStages := make(map[string]bool)
	for _, s := range stages.AvailableStages {
		knownStages[s.ID] = true
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		switch arg {
		case "-h", "--help", "help":
			cfg.showHelp = true
			return cfg
		case "-v", "--version", "version":
			cfg.showVersion = true
			return cfg
		case "-i", "--interactive", "repl":
			cfg.interactive = true
		case "--once":
			cfg.once = true
		case "--dangerously-skip-permissions":
			cfg.skipPermissions = true
		case "-p", "--print", "--prompt":
			if i+1 < len(args) {
				i++
				cfg.prompt = args[i]
				cfg.once = true
			}
		case "-s", "--speed":
			if i+1 < len(args) {
				i++
				cfg.speed = strings.ToLower(args[i])
			}
		case "-m", "--model":
			if i+1 < len(args) {
				i++
				cfg.model = args[i]
			}
		case "--stages":
			if i+1 < len(args) {
				i++
				cfg.includedStages = append(cfg.includedStages, args[i])
			}
		case "-e", "--exclude":
			// Consume one or more stage names following --exclude (like install-nothing --exclude cloud xorg)
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				cfg.excludedStages = append(cfg.excludedStages, args[i])
			}
		default:
			if strings.HasPrefix(arg, "--exclude=") {
				val := strings.TrimPrefix(arg, "--exclude=")
				cfg.excludedStages = append(cfg.excludedStages, val)
			} else if strings.HasPrefix(arg, "--stages=") {
				val := strings.TrimPrefix(arg, "--stages=")
				cfg.includedStages = append(cfg.includedStages, val)
			} else if strings.HasPrefix(arg, "--speed=") {
				cfg.speed = strings.ToLower(strings.TrimPrefix(arg, "--speed="))
			} else if strings.HasPrefix(arg, "--model=") {
				cfg.model = strings.TrimPrefix(arg, "--model=")
			} else if !strings.HasPrefix(arg, "-") {
				// Positional argument: check if it's a known stage name
				lower := strings.ToLower(arg)
				if knownStages[lower] {
					cfg.includedStages = append(cfg.includedStages, lower)
				} else {
					// Treat as custom prompt if not a stage
					if cfg.prompt == "" {
						cfg.prompt = arg
					} else {
						cfg.prompt += " " + arg
					}
				}
			}
		}
		i++
	}

	return cfg
}

func printHelp() {
	fmt.Printf("%s%srun-nothing%s (v%s) — A terminal application that simulates running Claude Code. It doesn't run anything.\n\n",
		ui.Bold, ui.Coral, ui.Reset, Version)
	fmt.Println("USAGE:")
	fmt.Println("  run-nothing [OPTIONS] [STAGES...]")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Run endless autonomous Claude Code simulation")
	fmt.Println("  run-nothing")
	fmt.Println()
	fmt.Println("  # Run specific stages only")
	fmt.Println("  run-nothing think codegen test")
	fmt.Println()
	fmt.Println("  # Exclude specific stages from simulation")
	fmt.Println("  run-nothing --exclude deploy review")
	fmt.Println()
	fmt.Println("  # Start interactive Claude Code REPL")
	fmt.Println("  run-nothing --interactive")
	fmt.Println()
	fmt.Println("  # Run a single custom prompt and exit")
	fmt.Println("  run-nothing -p \"Fix memory leak in websocket pool\"")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -i, --interactive                 Start interactive REPL mode")
	fmt.Println("  -p, --print <prompt>              Run a single prompt simulation and exit")
	fmt.Println("  -e, --exclude <stages...>         Exclude specific stages (e.g., --exclude deploy git)")
	fmt.Println("      --stages <stages...>          Run only specific stages")
	fmt.Println("  -s, --speed <mode>                Playback speed: fast, normal, cinematic, instant (default: normal)")
	fmt.Println("  -m, --model <model>               Model to display (default: claude-sonnet-4-6)")
	fmt.Println("      --once                        Run a single task cycle and exit")
	fmt.Println("      --dangerously-skip-permissions  Display simulated permission bypass warning")
	fmt.Println("  -h, --help                        Show this help message")
	fmt.Println("  -v, --version                     Show version number")

	stages.PrintAvailableStages()
}

func main() {
	cfg := parseArgs(os.Args[1:])

	if cfg.showHelp {
		printHelp()
		return
	}

	if cfg.showVersion {
		fmt.Printf("%s (run-nothing / Claude Code Simulator)\n", Version)
		return
	}

	// Setup graceful Ctrl+C handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
		fmt.Printf("\r\033[2K\n%sInterrupted by user (Ctrl+C). Exiting Claude Code simulation.%s\n", ui.Coral+ui.Bold, ui.Reset)
		fmt.Printf("%sActual tokens used: 0 · Actual cost: $0.00 · Files modified: 0%s\n", ui.Dim, ui.Reset)
		os.Exit(0)
	}()

	filter := stages.NewFilterConfig(cfg.includedStages, cfg.excludedStages)
	opts := sim.EngineOptions{
		Speed:           cfg.speed,
		Model:           cfg.model,
		Filter:          filter,
		InteractiveMode: cfg.interactive,
	}

	// Print Claude Code banner
	ui.PrintBanner(Version, cfg.model)

	if cfg.skipPermissions {
		fmt.Printf("%s⚠ WARNING: --dangerously-skip-permissions enabled. Bypassing all simulated safety checks.%s\n\n",
			ui.Yellow+ui.Bold, ui.Reset)
	}

	// Interactive REPL mode
	if cfg.interactive {
		repl.RunInteractiveLoop(ctx, opts, Version)
		return
	}

	// Single custom prompt mode
	if cfg.prompt != "" {
		sc := sim.BuildCustomScenario(cfg.prompt)
		_ = sim.RunScenario(ctx, sc, opts)
		return
	}

	// Continuous autonomous mode (like install-nothing)
	fmt.Printf("%sAutonomous Agent Mode active · Press Ctrl+C to stop%s\n\n", ui.Dim, ui.Reset)

	taskIdx := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		sc := sim.Scenarios[taskIdx%len(sim.Scenarios)]
		err := sim.RunScenario(ctx, sc, opts)
		if err != nil {
			return
		}

		if cfg.once {
			return
		}

		taskIdx++
		opts.Sleep(1000 * time.Millisecond)
		fmt.Printf("%s────────────────────────────────────────────────────────────%s\n\n", ui.DarkGray, ui.Reset)
	}
}
