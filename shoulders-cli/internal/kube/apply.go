package kube

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func Apply(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, obj *unstructured.Unstructured) error {
	var resource dynamic.ResourceInterface = client.Resource(gvr)
	if namespace != "" {
		resource = client.Resource(gvr).Namespace(namespace)
	}

	current, err := resource.Get(ctx, obj.GetName(), metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = resource.Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}
	obj.SetResourceVersion(current.GetResourceVersion())
	_, err = resource.Update(ctx, obj, metav1.UpdateOptions{})
	return err
}
