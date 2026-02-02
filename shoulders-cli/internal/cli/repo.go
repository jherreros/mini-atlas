package cli

import (
	"errors"
	"os"
	"path/filepath"
)

func FindRepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	current := cwd
	for {
		candidate := filepath.Join(current, "1-cluster", "create-cluster.sh")
		if _, err := os.Stat(candidate); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", errors.New("could not locate repo root (missing 1-cluster/create-cluster.sh)")
		}
		current = parent
	}
}
