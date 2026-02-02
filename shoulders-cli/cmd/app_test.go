package cmd

import "testing"

func TestParseImageTag(t *testing.T) {
	image, tag := parseImageTag("nginx:1.26", "")
	if image != "nginx" || tag != "1.26" {
		t.Fatalf("expected nginx:1.26, got %s:%s", image, tag)
	}

	image, tag = parseImageTag("nginx", "")
	if tag != "latest" {
		t.Fatalf("expected latest, got %s", tag)
	}

	image, tag = parseImageTag("nginx", "custom")
	if tag != "custom" {
		t.Fatalf("expected custom, got %s", tag)
	}
}
