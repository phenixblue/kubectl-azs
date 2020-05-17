package kube

import (
	"k8s.io/client-go/kubernetes"

	// Import all auth client plugins

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateKubeClient creates a new kubernetes client interface
func CreateKubeClient(kubeconfig string, configContext string) (kubernetes.Interface, error) {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	loadingRules.ExplicitPath = kubeconfig
	configOverrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults, CurrentContext: configContext}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	clientset, _ := kubernetes.NewForConfig(config)

	return clientset, err

}
