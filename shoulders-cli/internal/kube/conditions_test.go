package kube

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestHasCondition(t *testing.T) {
	tests := []struct {
		name           string
		obj            unstructured.Unstructured
		condType       string
		condStatus     string
		expectedResult bool
	}{
		{
			name: "ConditionExistsAndMatches",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Ready",
								"status": "True",
							},
						},
					},
				},
			},
			condType:       "Ready",
			condStatus:     "True",
			expectedResult: true,
		},
		{
			name: "ConditionExistsButMismatch",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Ready",
								"status": "False",
							},
						},
					},
				},
			},
			condType:       "Ready",
			condStatus:     "True",
			expectedResult: false,
		},
		{
			name: "ConditionDoesNotExist",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Other",
								"status": "True",
							},
						},
					},
				},
			},
			condType:       "Ready",
			condStatus:     "True",
			expectedResult: false,
		},
		{
			name: "NoStatus",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			condType:       "Ready",
			condStatus:     "True",
			expectedResult: false,
		},
		{
			name: "NoConditions",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			condType:       "Ready",
			condStatus:     "True",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := HasCondition(tt.obj, tt.condType, tt.condStatus)
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
