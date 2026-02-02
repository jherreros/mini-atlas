package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jherreros/shoulders/shoulders-cli/internal/bootstrap"
	"github.com/jherreros/shoulders/shoulders-cli/internal/cli"
	"github.com/jherreros/shoulders/shoulders-cli/internal/flux"
	"github.com/jherreros/shoulders/shoulders-cli/internal/kube"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create the local cluster and install platform addons",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := cli.FindRepoRoot()
		if err != nil {
			return err
		}

		configPath := filepath.Join(repoRoot, "1-cluster", "kind-config.yaml")
		if err := bootstrap.EnsureKindCluster(bootstrap.DefaultClusterName, configPath); err != nil {
			return fmt.Errorf("failed to create kind cluster: %w", err)
		}
		if err := bootstrap.EnsureCilium(kubeconfig); err != nil {
			return fmt.Errorf("failed to install cilium: %w", err)
		}
		if err := bootstrap.EnsureFlux(context.Background(), kubeconfig, repoRoot); err != nil {
			return fmt.Errorf("failed to install flux: %w", err)
		}

		spinner, err := pterm.DefaultSpinner.Start("Waiting for Flux to reconcile")
		if err != nil {
			return err
		}
		if err := waitForFlux(spinner); err != nil {
			spinner.Fail("Flux reconciliation failed")
			return err
		}
		spinner.Success("Flux reconciliation complete")
		return nil
	},
}

func waitForFlux(spinner *pterm.SpinnerPrinter) error {
	ctx := context.Background()
	client, err := kube.NewDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-ticker.C:
			ready, pending, err := flux.AllKustomizationsReady(ctx, client, "flux-system")
			if err != nil {
				return err
			}
			if ready {
				return nil
			}
			spinner.UpdateText(fmt.Sprintf("Waiting for Flux: %v", pending))
		case <-timeout:
			return fmt.Errorf("timed out waiting for Flux Kustomizations")
		}
	}
}
