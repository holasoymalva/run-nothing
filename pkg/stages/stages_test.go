package stages

import "testing"

func TestFilterConfig(t *testing.T) {
	// Test default all enabled
	filter := NewFilterConfig(nil, nil)
	if !filter.IsEnabled("scout") || !filter.IsEnabled("think") {
		t.Errorf("Expected all stages enabled by default")
	}

	// Test inclusion
	includeFilter := NewFilterConfig([]string{"think", "codegen"}, nil)
	if !includeFilter.IsEnabled("think") {
		t.Errorf("Expected 'think' to be enabled")
	}
	if !includeFilter.IsEnabled("codegen") {
		t.Errorf("Expected 'codegen' to be enabled")
	}
	if includeFilter.IsEnabled("deploy") {
		t.Errorf("Expected 'deploy' to be disabled")
	}

	// Test exclusion
	excludeFilter := NewFilterConfig(nil, []string{"deploy", "review"})
	if excludeFilter.IsEnabled("deploy") {
		t.Errorf("Expected 'deploy' to be excluded")
	}
	if excludeFilter.IsEnabled("review") {
		t.Errorf("Expected 'review' to be excluded")
	}
	if !excludeFilter.IsEnabled("think") {
		t.Errorf("Expected 'think' to remain enabled")
	}

	// Test comma separated format
	commaFilter := NewFilterConfig(nil, []string{"git,deploy"})
	if commaFilter.IsEnabled("git") || commaFilter.IsEnabled("deploy") {
		t.Errorf("Expected comma-separated excluded stages to be disabled")
	}
	if !commaFilter.IsEnabled("test") {
		t.Errorf("Expected 'test' to remain enabled")
	}
}
