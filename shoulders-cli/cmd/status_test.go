package cmd

import "testing"

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name   string
		ready  bool
		issues []string
		want   string
	}{
		{"Ready", true, nil, "Healthy"},
		{"ReadyWithIssues", true, []string{"warn"}, "Healthy"}, // Should probably not happen logic-wise but good to assert behavior
		{"NotReadyNoIssues", false, nil, "Unhealthy"},
		{"NotReadyWithIssues", false, []string{"err1"}, "Unhealthy (err1)"},
		{"NotReadyMultipleIssues", false, []string{"err1", "err2"}, "Unhealthy (err1, err2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatStatus(tt.ready, tt.issues); got != tt.want {
				t.Errorf("formatStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolToText(t *testing.T) {
	if got := boolToText(true); got != "Ready" {
		t.Errorf("boolToText(true) = %v, want Ready", got)
	}
	if got := boolToText(false); got != "Not Ready" {
		t.Errorf("boolToText(false) = %v, want Not Ready", got)
	}
}
