package kube

import (
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GetNodes gets a list of Kubernetes nodes filtered by the target AZ label
func GetNodes(client kubernetes.Interface, label string, labelChanged bool) (*v1.NodeList, string, error) {

	var nodes *v1.NodeList
	origLabel := label

	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(nodes.Items) < 1 {

		if !labelChanged {

			label = strings.Replace(label, "failure-domain", "failure-domain.beta", 1)

			fmt.Printf("Beta Label -> %v\n", label)

			nodes, err = client.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: label})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}

		if len(nodes.Items) < 1 {

			fmt.Printf("No nodes with target AZ label (%q) found\n", origLabel)
			os.Exit(1)
		}

	}

	return nodes, label, err

}
