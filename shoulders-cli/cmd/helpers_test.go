package cmd

import (
	"testing"

	"github.com/jherreros/shoulders/shoulders-cli/internal/config"
)

func TestCurrentNamespaceOverride(t *testing.T) {
	originalOverride := namespaceOverride
	originalConfig := currentConfig
	defer func() {
		namespaceOverride = originalOverride
		currentConfig = originalConfig
	}()

	namespaceOverride = "override"
	currentConfig = &config.Config{CurrentWorkspace: "team-a"}

	ns, err := currentNamespace()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns != "override" {
		t.Fatalf("expected override namespace, got %s", ns)
	}
}

func TestCurrentNamespaceFromConfig(t *testing.T) {
	originalOverride := namespaceOverride
	originalConfig := currentConfig
	defer func() {
		namespaceOverride = originalOverride
		currentConfig = originalConfig
	}()

	namespaceOverride = ""
	currentConfig = &config.Config{CurrentWorkspace: "team-b"}

	ns, err := currentNamespace()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns != "team-b" {
		t.Fatalf("expected team-b namespace, got %s", ns)
	}
}

func TestCurrentNamespaceMissing(t *testing.T) {
	originalOverride := namespaceOverride
	originalConfig := currentConfig
	defer func() {
		namespaceOverride = originalOverride
		currentConfig = originalConfig
	}()

	namespaceOverride = ""
	currentConfig = nil

	_, err := currentNamespace()
	if err == nil {
		t.Fatalf("expected error when namespace missing")
	}
}
