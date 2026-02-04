package bootstrap

import (
	"fmt"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

const DefaultClusterName = "shoulders"

func EnsureKindCluster(name, configPath string) error {
	provider, err := newProvider()
	if err != nil {
		return err
	}
	clusters, err := provider.List()
	if err != nil {
		return err
	}
	for _, existing := range clusters {
		if existing == name {
			return nil
		}
	}
	return provider.Create(
		name,
		cluster.CreateWithConfigFile(configPath),
		cluster.CreateWithWaitForReady(5*time.Minute),
	)
}

func DeleteKindCluster(name string) error {
	provider, err := newProvider()
	if err != nil {
		return err
	}
	return provider.Delete(name, "")
}

func ListClusters() ([]string, error) {
	provider, err := newProvider()
	if err != nil {
		return nil, err
	}
	return provider.List()
}

func newProvider() (*cluster.Provider, error) {
	opt, err := cluster.DetectNodeProvider()
	if err != nil {
		return nil, fmt.Errorf("no container runtime detected; please install Docker or Podman: %w", err)
	}
	return cluster.NewProvider(opt, cluster.ProviderWithLogger(log.NoopLogger{})), nil
}
