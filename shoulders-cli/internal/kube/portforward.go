package kube

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const portForwardTimeout = 30 * time.Second

func PortForwardService(ctx context.Context, kubeconfigPath, namespace, serviceName string, localPort, remotePort int) (chan struct{}, chan struct{}, error) {
	config, err := NewRestConfig(kubeconfigPath)
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	service, err := clientset.CoreV1().Services(namespace).Get(ctx, serviceName, getOptions())
	if err != nil {
		return nil, nil, err
	}
	podName, err := findServicePod(ctx, clientset, namespace, service.Spec.Selector)
	if err != nil {
		return nil, nil, err
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return nil, nil, err
	}

	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})
	ports := []string{fmt.Sprintf("%d:%d", localPort, remotePort)}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())

	forwarder, err := portforward.NewOnAddresses(dialer, []string{"127.0.0.1"}, ports, stopCh, readyCh, os.Stdout, os.Stderr)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		_ = forwarder.ForwardPorts()
	}()

	select {
	case <-readyCh:
		return stopCh, readyCh, nil
	case <-time.After(portForwardTimeout):
		close(stopCh)
		return nil, nil, fmt.Errorf("port-forward timeout for service %s/%s", namespace, serviceName)
	case <-ctx.Done():
		close(stopCh)
		return nil, nil, ctx.Err()
	}
}

func findServicePod(ctx context.Context, clientset *kubernetes.Clientset, namespace string, selector map[string]string) (string, error) {
	if len(selector) == 0 {
		return "", errors.New("service has no selector")
	}
	requirements := []labels.Requirement{}
	for key, value := range selector {
		req, err := labels.NewRequirement(key, selection.Equals, []string{value})
		if err != nil {
			return "", err
		}
		requirements = append(requirements, *req)
	}
	selectorObj := labels.NewSelector().Add(requirements...)
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, listOptions(selectorObj.String()))
	if err != nil {
		return "", err
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			return pod.Name, nil
		}
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pods found for service selector %s", selectorString(selector))
	}
	return pods.Items[0].Name, nil
}

func listOptions(labelSelector string) metav1.ListOptions {
	return metav1.ListOptions{LabelSelector: labelSelector}
}

func selectorString(selector map[string]string) string {
	pairs := make([]string, 0, len(selector))
	for key, value := range selector {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(pairs, ",")
}

func getOptions() metav1.GetOptions {
	return metav1.GetOptions{}
}
