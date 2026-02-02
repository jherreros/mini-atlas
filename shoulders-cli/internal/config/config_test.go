package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	viper.Reset()

	cfg := &Config{CurrentWorkspace: "team-a"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.CurrentWorkspace != "team-a" {
		t.Fatalf("expected workspace team-a, got %s", loaded.CurrentWorkspace)
	}

	configPath, err := Path()
	if err != nil {
		t.Fatalf("path failed: %v", err)
	}
	expected := filepath.Join(tmpDir, ".shoulders", "config.yaml")
	if configPath != expected {
		t.Fatalf("expected config path %s, got %s", expected, configPath)
	}
}
