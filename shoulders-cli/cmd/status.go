package cmd

import (
	"context"
	"fmt"

	"github.com/jherreros/shoulders/shoulders-cli/internal/crossplane"
	"github.com/jherreros/shoulders/shoulders-cli/internal/flux"
	"github.com/jherreros/shoulders/shoulders-cli/internal/kube"
	"github.com/jherreros/shoulders/shoulders-cli/internal/output"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type statusSummary struct {
	NodesReady   bool     `json:"nodesReady" yaml:"nodesReady"`
	FluxReady    bool     `json:"fluxReady" yaml:"fluxReady"`
	FluxPending  []string `json:"fluxPending" yaml:"fluxPending"`
	ProvidersOK  bool     `json:"providersHealthy" yaml:"providersHealthy"`
	ProvidersBad []string `json:"providersUnhealthy" yaml:"providersUnhealthy"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cluster and platform status",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, err := outputOption()
		if err != nil {
			return err
		}

		ctx := context.Background()
		clientset, err := kube.NewClientset(kubeconfig)
		if err != nil {
			return err
		}
		dynamicClient, err := kube.NewDynamicClient(kubeconfig)
		if err != nil {
			return err
		}

		nodesReady := true
		nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
		if err != nil {
			return err
		}
		for _, node := range nodes.Items {
			ready := false
			for _, cond := range node.Status.Conditions {
				if cond.Type == "Ready" && cond.Status == "True" {
					ready = true
				}
			}
			if !ready {
				nodesReady = false
			}
		}

		fluxReady, pending, err := flux.AllKustomizationsReady(ctx, dynamicClient, "flux-system")
		if err != nil {
			return err
		}
		providersOK, unhealthy, err := crossplane.AllProvidersHealthy(ctx, dynamicClient)
		if err != nil {
			return err
		}

		summary := statusSummary{
			NodesReady:   nodesReady,
			FluxReady:    fluxReady,
			FluxPending:  pending,
			ProvidersOK:  providersOK,
			ProvidersBad: unhealthy,
		}

		if format == output.Table {
			rows := [][]string{
				{"Nodes Ready", fmt.Sprintf("%t", summary.NodesReady)},
				{"Flux Ready", fmt.Sprintf("%t", summary.FluxReady)},
				{"Flux Pending", fmt.Sprintf("%v", summary.FluxPending)},
				{"Crossplane Providers Healthy", fmt.Sprintf("%t", summary.ProvidersOK)},
				{"Unhealthy Providers", fmt.Sprintf("%v", summary.ProvidersBad)},
			}
			return output.PrintTable([]string{"Check", "Result"}, rows)
		}

		payload, err := output.Render(summary, format)
		if err != nil {
			return err
		}
		fmt.Println(string(payload))
		return nil
	},
}
