package crossplane

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var providerGVR = schema.GroupVersionResource{
	Group:    "pkg.crossplane.io",
	Version:  "v1",
	Resource: "providers",
}

func ListProviders(ctx context.Context, client dynamic.Interface) ([]unstructured.Unstructured, error) {
	list, err := client.Resource(providerGVR).List(ctx, listOptions())
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func AllProvidersHealthy(ctx context.Context, client dynamic.Interface) (bool, []string, error) {
	providers, err := ListProviders(ctx, client)
	if err != nil {
		return false, nil, err
	}
	var unhealthy []string
	for _, provider := range providers {
		ready, _ := hasCondition(provider, "Healthy", "True")
		if !ready {
			unhealthy = append(unhealthy, provider.GetName())
		}
	}
	return len(unhealthy) == 0, unhealthy, nil
}

func listOptions() metav1.ListOptions {
	return metav1.ListOptions{}
}

func hasCondition(obj unstructured.Unstructured, conditionType, expectedStatus string) (bool, error) {
	status, ok := obj.Object["status"].(map[string]interface{})
	if !ok {
		return false, nil
	}
	conditions, ok := status["conditions"].([]interface{})
	if !ok {
		return false, nil
	}
	for _, cond := range conditions {
		condition, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}
		if condition["type"] == conditionType && condition["status"] == expectedStatus {
			return true, nil
		}
	}
	return false, nil
}
