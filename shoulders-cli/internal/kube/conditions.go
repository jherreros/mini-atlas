package kube

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// HasCondition checks if an unstructured object has a specific condition with a specific status.
func HasCondition(obj unstructured.Unstructured, conditionType, expectedStatus string) (bool, error) {
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
