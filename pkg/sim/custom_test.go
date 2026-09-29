package sim

import (
	"strings"
	"testing"
)

func TestBuildCustomScenario(t *testing.T) {
	prompt := "Refactor token bucket rate limiter to use atomic CAS"
	sc := BuildCustomScenario(prompt)

	if sc.Prompt != prompt {
		t.Errorf("Expected prompt to match, got %s", sc.Prompt)
	}

	if len(sc.Thoughts) == 0 {
		t.Errorf("Expected thoughts to be populated")
	}

	if len(sc.Diffs) == 0 {
		t.Errorf("Expected diffs to be populated")
	}

	if !strings.Contains(sc.TestCommand, "go test") {
		t.Errorf("Expected test command to run go test, got: %s", sc.TestCommand)
	}

	if !strings.Contains(sc.Conclusion, "statement coverage") {
		t.Errorf("Expected conclusion to report statement coverage")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Fix memory leak in pool", "fix-memory"},
		{"Refactor auth middleware", "refactor-auth"},
		{"Add JWT verification", "add-jwt"},
	}

	for _, tt := range tests {
		got := slugify(tt.input)
		if got != tt.expected {
			t.Errorf("slugify(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
