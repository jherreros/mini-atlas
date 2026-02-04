package flux

import (
	"context"
	"fmt"

	"github.com/jherreros/shoulders/shoulders-cli/internal/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var kustomizationGVR = schema.GroupVersionResource{
	Group:    "kustomize.toolkit.fluxcd.io",
	Version:  "v1",
	Resource: "kustomizations",
}

func ListKustomizations(ctx context.Context, client dynamic.Interface, namespace string) ([]unstructured.Unstructured, error) {
	resource := client.Resource(kustomizationGVR)
	var listResource dynamic.ResourceInterface = resource
	if namespace != "" {
		listResource = resource.Namespace(namespace)
	}
	list, err := listResource.List(ctx, listOptions())
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func AllKustomizationsReady(ctx context.Context, client dynamic.Interface, namespace string) (bool, []string, error) {
	items, err := ListKustomizations(ctx, client, namespace)
	if err != nil {
		return false, nil, err
	}
	var notReady []string
	for _, item := range items {
		name := item.GetName()
		ready, _ := kube.HasCondition(item, "Ready", "True")
		if !ready {
			notReady = append(notReady, name)
		}
	}
	return len(notReady) == 0, notReady, nil
}

func listOptions() metav1.ListOptions {
	return metav1.ListOptions{}
}

func KustomizationStatusSummary(ctx context.Context, client dynamic.Interface, namespace string) (string, error) {
	ready, notReady, err := AllKustomizationsReady(ctx, client, namespace)
	if err != nil {
		return "", err
	}
	if ready {
		return "All Kustomizations Ready", nil
	}
	return fmt.Sprintf("Not Ready: %v", notReady), nil
}
