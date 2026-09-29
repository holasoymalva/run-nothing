package ui

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Colors
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	// Claude branding coral/orange
	Coral    = "\033[38;5;208m"
	Purple   = "\033[38;5;141m"
	SoftCyan = "\033[38;5;80m"
	DarkGray = "\033[38;5;238m"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

// PrintBanner renders the authentic Claude Code welcome banner
func PrintBanner(version string, model string) {
	width := 64
	top    := "╭" + strings.Repeat("─", width-2) + "╮"
	bottom := "╰" + strings.Repeat("─", width-2) + "╯"

	line1 := fmt.Sprintf("   %s%sClaude Code%s %s(v%s)%s", Bold, Coral, Reset, Dim, version, Reset)
	line2 := fmt.Sprintf("   %sOpus 4.6%s · %sSonnet 4.6%s · %sHaiku 4.5%s", Dim, Reset, Cyan, Reset, Dim, Reset)
	line3 := fmt.Sprintf("   Welcome to %sClaude Code%s research preview! %s(simulated)%s", Coral, Reset, Dim, Reset)
	line4 := fmt.Sprintf("   Active model: %s%s%s · Type %s/help%s for commands", Bold, model, Reset, Cyan, Reset)

	fmt.Println()
	fmt.Println(Coral + top + Reset)
	printBoxLine("", width)
	printBoxLine(line1, width)
	printBoxLine(line2, width)
	printBoxLine("", width)
	printBoxLine(line3, width)
	printBoxLine(line4, width)
	printBoxLine("", width)
	fmt.Println(Coral + bottom + Reset)
	fmt.Println()
}

func printBoxLine(content string, width int) {
	visibleLen := utf8.RuneCountInString(stripANSI(content))
	padding := width - 2 - visibleLen
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s│%s%s%s%s│%s\n", Coral, Reset, content, strings.Repeat(" ", padding), Coral, Reset)
}

// Spinner frames
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner represents an animated terminal spinner
type Spinner struct {
	label   string
	stopCh  chan struct{}
	doneCh  chan struct{}
	speedMs time.Duration
}

// StartSpinner starts a background animated spinner
func StartSpinner(label string) *Spinner {
	s := &Spinner{
		label:   label,
		stopCh:  make(chan struct{}),
		doneCh:  make(chan struct{}),
		speedMs: 80 * time.Millisecond,
	}

	go func() {
		defer close(s.doneCh)
		i := 0
		for {
			select {
			case <-s.stopCh:
				fmt.Print("\r\033[2K") // Clear the line
				return
			default:
				frame := spinnerFrames[i%len(spinnerFrames)]
				fmt.Printf("\r%s%s%s %s%s%s", Coral, frame, Reset, Dim, s.label, Reset)
				i++
				time.Sleep(s.speedMs)
			}
		}
	}()

	return s
}

// UpdateLabel changes the spinner's text dynamically
func (s *Spinner) UpdateLabel(newLabel string) {
	s.label = newLabel
}

// Stop ends the spinner and cleans up line
func (s *Spinner) Stop(finalSymbol string, finalMessage string) {
	close(s.stopCh)
	<-s.doneCh
	if finalSymbol != "" && finalMessage != "" {
		fmt.Printf("\r\033[2K%s %s\n", finalSymbol, finalMessage)
	} else {
		fmt.Print("\r\033[2K")
	}
}

// ProgressBar renders an authentic install-nothing style progress bar
func RenderProgressBar(current, total int, width int, label string) {
	percent := float64(current) / float64(total)
	filledLen := int(float64(width) * percent)
	if filledLen > width {
		filledLen = width
	}
	emptyLen := width - filledLen

	bar := strings.Repeat("█", filledLen) + strings.Repeat("░", emptyLen)
	pctStr := fmt.Sprintf("%3.0f%%", percent*100)

	fmt.Printf("\r%s%s%s [%s%s%s] %s%s%s %s%s%s",
		Cyan, label, Reset,
		Green, bar, Reset,
		Bold, pctStr, Reset,
		Dim, fmt.Sprintf("(%d/%d)", current, total), Reset)

	if current == total {
		fmt.Println()
	}
}

// Typewriter prints text with human/AI typing cadence
func Typewriter(text string, charDelay time.Duration) {
	for _, r := range text {
		fmt.Print(string(r))
		if r == '\n' {
			time.Sleep(charDelay * 2)
		} else if r == '.' || r == ',' || r == ':' {
			time.Sleep(charDelay * 3)
		} else {
			jitter := time.Duration(rand.Intn(4)-2) * time.Millisecond
			sleepTime := charDelay + jitter
			if sleepTime < 1*time.Millisecond {
				sleepTime = 1 * time.Millisecond
			}
			time.Sleep(sleepTime)
		}
	}
	fmt.Println()
}

// PrintUnifiedDiff renders colorized git diff output
func PrintUnifiedDiff(filename string, diffLines []string) {
	fmt.Println()
	fmt.Printf("%s%s--- a/%s%s\n", Bold, Red, filename, Reset)
	fmt.Printf("%s%s+++ b/%s%s\n", Bold, Green, filename, Reset)

	for _, line := range diffLines {
		if strings.HasPrefix(line, "@@") {
			fmt.Printf("%s%s%s\n", Cyan, line, Reset)
		} else if strings.HasPrefix(line, "+") {
			fmt.Printf("%s%s%s\n", Green, line, Reset)
		} else if strings.HasPrefix(line, "-") {
			fmt.Printf("%s%s%s\n", Red, line, Reset)
		} else {
			fmt.Printf("%s%s%s\n", Gray, line, Reset)
		}
		time.Sleep(25 * time.Millisecond)
	}
	fmt.Println()
}

// PrintTurnSummary renders Claude Code's bottom turn stats
func PrintTurnSummary(cost float64, inTokens, outTokens int, duration time.Duration) {
	fmt.Printf("%s● Cost: $%0.4f · Tokens: %d in / %d out · Duration: %0.1fs%s\n",
		Dim, cost, inTokens, outTokens, duration.Seconds(), Reset)
	fmt.Println()
}

// PrintToolCall renders a simulated tool invocation
func PrintToolCall(toolName string, args string) {
	fmt.Printf("%s●%s %s%s%s%s(%s)%s\n",
		Coral, Reset, Bold, toolName, Reset,
		Dim, args, Reset)
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// IsTerminal returns whether stdout is a terminal
func IsTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
