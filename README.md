# run-nothing

A terminal application that simulates running Claude Code. It doesn't actually run anything.

```
╭──────────────────────────────────────────────────────────────╮
│                                                              │
│   Claude Code (v2.1.139)                                     │
│   Opus 4.6 · Sonnet 4.6 · Haiku 4.5                          │
│                                                              │
│   Welcome to Claude Code research preview! (simulated)       │
│   Active model: claude-sonnet-4-6 · Type /help for commands  │
│                                                              │
╰──────────────────────────────────────────────────────────────╯
```

Look extraordinarily productive while doing absolutely nothing. Simulates deep chain-of-thought reasoning, codebase indexing, colorized unified git diffs, test runners, multi-agent reviews, and canary deployments. Zero Anthropic API tokens used, $0.00 spent, zero local files touched.

## Installation

### Download binary

Grab the latest binary for your platform from [Releases](https://github.com/holasoymalva/claude-nothing/releases)

```bash
chmod +x run-nothing-*
./run-nothing-linux-x86_64
```

### Homebrew

```bash
brew install run-nothing
```

### Build from source

```bash
go run .
# or build a binary
make build
./run-nothing
```

Press `Ctrl+C` to stop.

### Pick what to simulate

By default we simulate everything. But you can change this behavior.

```bash
# Simulate specific stages
./run-nothing think codegen
```

Or pick what not to simulate:

```bash
# Exclude specific stages from simulation
./run-nothing --exclude deploy review
```

See available stages:
```bash
./run-nothing --help
```

Available stages:
- `scout` - Codebase exploration (indexing, ripgrep, AST parsing, reading files)
- `think` - Deep reasoning (animated CoT spinner with technical thoughts)
- `codegen` - Code synthesis (colorized unified git diffs)
- `test` - Test runner & verification (passing unit tests, race detectors)
- `review` - Multi-agent code review (automated UltraReview security audit)
- `git` - Version control (atomic commits, conventional messages, branches)
- `deploy` - CI/CD & canary rollout (pod synchronization progress bar)

### Interactive Mode

You can also run it as an interactive Claude Code REPL:

```bash
./run-nothing --interactive
```

Chat with simulated Claude, run slash commands (`/help`, `/cost`, `/model`, `/doctor`, `/compact`, `/clear`, `/exit`), and prompt it to "solve" any imaginary problem!

### Single Prompt Mode

Simulate answering a specific task directly:

```bash
./run-nothing -p "Refactor token bucket rate limiter to use atomic CAS"
```

### Playback Speed

Control the simulation pacing:

```bash
# Options: fast, normal, cinematic, instant
./run-nothing --speed cinematic
```

## Docker

Build
```bash
docker build -t run-nothing .
```

Run
```bash
docker run -it --rm --init run-nothing
```

## Nix

Install
```bash
nix profile install github:holasoymalva/claude-nothing
```

Run
```bash
nix run github:holasoymalva/claude-nothing
```

## License

Do whatever you want with it. Well, except for movies. If you use this in a movie, credit me or something.
