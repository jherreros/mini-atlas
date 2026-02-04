package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jherreros/shoulders/shoulders-cli/internal/crossplane"
	"github.com/jherreros/shoulders/shoulders-cli/internal/flux"
	"github.com/jherreros/shoulders/shoulders-cli/internal/kube"
	"github.com/jherreros/shoulders/shoulders-cli/internal/output"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type statusSummary struct {
	K8sVersion   string   `json:"k8sVersion" yaml:"k8sVersion"`
	NodesReady   bool     `json:"nodesReady" yaml:"nodesReady"`
	NodeCount    int      `json:"nodeCount" yaml:"nodeCount"`
	FluxReady    bool     `json:"fluxReady" yaml:"fluxReady"`
	FluxBroken   []string `json:"fluxBroken" yaml:"fluxBroken"`
	XPlaneReady  bool     `json:"crossplaneReady" yaml:"crossplaneReady"`
	XPlaneBroken []string `json:"crossplaneBroken" yaml:"crossplaneBroken"`
	GatewayReady bool     `json:"gatewayReady" yaml:"gatewayReady"`
	GatewayAddr  string   `json:"gatewayAddress" yaml:"gatewayAddress"`
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

		// 1. Cluster Info
		versionInfo, err := clientset.Discovery().ServerVersion()
		k8sVersion := "unknown"
		if err == nil {
			k8sVersion = versionInfo.GitVersion
		}

		// 2. Nodes
		nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
		if err != nil {
			return err
		}
		nodesReady := true
		for _, node := range nodes.Items {
			if !isNodeReady(node) {
				nodesReady = false
			}
		}

		// 3. Flux
		fluxReady, fluxPending, err := flux.AllKustomizationsReady(ctx, dynamicClient, "flux-system")
		if err != nil {
			// Don't fail completely if CRDs are missing, just report not ready
			fluxReady = false
			fluxPending = []string{err.Error()}
		}

		// 4. Crossplane
		xpReady, xpUnhealthy, err := crossplane.AllProvidersHealthy(ctx, dynamicClient)
		if err != nil {
			xpReady = false
			xpUnhealthy = []string{err.Error()}
		}

		// 5. Gateway
		gwReady := false
		gwAddr := "Pending"
		gvrGW := schema.GroupVersionResource{Group: "gateway.networking.k8s.io", Version: "v1", Resource: "gateways"}
		gwList, err := dynamicClient.Resource(gvrGW).Namespace("gateway").List(ctx, v1.ListOptions{})
		if err == nil && len(gwList.Items) > 0 {
			// Check the first gateway found (usually 'main')
			gw := gwList.Items[0]
			gwReady = true // Simplified check

			// Try to find address in status
			if status, ok := gw.Object["status"].(map[string]interface{}); ok {
				if addrs, ok := status["addresses"].([]interface{}); ok && len(addrs) > 0 {
					if addrMap, ok := addrs[0].(map[string]interface{}); ok {
						gwAddr = fmt.Sprintf("%v", addrMap["value"])
					}
				}
			}
		}

		summary := statusSummary{
			K8sVersion:   k8sVersion,
			NodesReady:   nodesReady,
			NodeCount:    len(nodes.Items),
			FluxReady:    fluxReady,
			FluxBroken:   fluxPending,
			XPlaneReady:  xpReady,
			XPlaneBroken: xpUnhealthy,
			GatewayReady: gwReady,
			GatewayAddr:  gwAddr,
		}

		if format == output.Table {
			printStatusTable(summary)
			return nil
		}

		payload, err := output.Render(summary, format)
		if err != nil {
			return err
		}
		fmt.Println(string(payload))
		return nil
	},
}

func isNodeReady(node corev1.Node) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func printStatusTable(s statusSummary) {
	fmt.Println("CLUSTER")
	fmt.Printf("  Version: %s\n", s.K8sVersion)
	fmt.Printf("  Nodes:   %d (All Ready: %t)\n", s.NodeCount, s.NodesReady)
	fmt.Println()

	fmt.Println("PLATFORM")
	fmt.Printf("  Flux CD:    %s\n", formatStatus(s.FluxReady, s.FluxBroken))
	fmt.Printf("  Crossplane: %s\n", formatStatus(s.XPlaneReady, s.XPlaneBroken))
	fmt.Printf("  Gateway:    %s (%s)\n", boolToText(s.GatewayReady), s.GatewayAddr)
}

func formatStatus(ready bool, issues []string) string {
	if ready {
		return "Healthy"
	}
	if len(issues) > 0 {
		return fmt.Sprintf("Unhealthy (%s)", strings.Join(issues, ", "))
	}
	return "Unhealthy"
}

func boolToText(b bool) string {
	if b {
		return "Ready"
	}
	return "Not Ready"
}
