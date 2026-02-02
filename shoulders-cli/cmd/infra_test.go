package cmd

import "testing"

func TestParseConfig(t *testing.T) {
	config, err := parseConfig([]string{"cleanup.policy=compact", "retention.ms=60000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config["cleanup.policy"] != "compact" {
		t.Fatalf("expected cleanup.policy=compact")
	}
}

func TestParseConfigInvalid(t *testing.T) {
	_, err := parseConfig([]string{"invalid"})
	if err == nil {
		t.Fatalf("expected error for invalid entry")
	}
}
