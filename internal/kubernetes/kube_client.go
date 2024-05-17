package kubernetes

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubeClient(targetNamespace string) (*kubernetes.Clientset, error) {
	namespace, err := getNamespaceFromServiceAccount()
	if err != nil {
		return nil, fmt.Errorf("error getting namespace from service account: %v", err)
	}

	kubeconfigContent, err := getKubeconfigFromConfigMap(namespace, "kubeconfig", targetNamespace)
	if err != nil {
		return nil, fmt.Errorf("error reading kubeconfig from ConfigMap: %v", err)
	}

	// Create Kubernetes client
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfigContent))
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return clientset, nil
}

func getKubeconfigFromConfigMap(namespace, name, key string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", fmt.Errorf("error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error getting ConfigMap: %v", err)
	}

	kubeconfig, ok := configMap.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in ConfigMap %s", key, name)
	}

	return kubeconfig, nil
}

func getNamespaceFromServiceAccount() (string, error) {
	namespaceFile := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespace, err := os.ReadFile(namespaceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read namespace file: %v", err)
	}
	return string(namespace), nil
}
