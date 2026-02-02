package kube

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func NewDynamicClient(kubeconfig string) (dynamic.Interface, error) {
	config, err := NewRestConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func NewClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := NewRestConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
