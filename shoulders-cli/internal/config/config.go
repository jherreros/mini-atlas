package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	CurrentWorkspace string `mapstructure:"current_workspace" yaml:"current_workspace" json:"current_workspace"`
}

func DefaultConfig() *Config {
	return &Config{}
}

func Load() (*Config, error) {
	cfg := DefaultConfig()
	configPath, err := Path()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err == nil {
		if err := viper.Unmarshal(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	configPath, err := Path()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	viper.SetConfigFile(configPath)
	viper.Set("current_workspace", cfg.CurrentWorkspace)
	return viper.WriteConfigAs(configPath)
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".shoulders", "config.yaml"), nil
}
