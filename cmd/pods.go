// Copyright © 2019 Joe Searcy <joe@twr.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"kubectl-azs/pkg/kube"
	"kubectl-azs/pkg/printers"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// podObj is an instance of a Kubernetes Pod Object
type podObj struct {
	name      string
	namespace string
	node      string
	az        string
}

// nodeObj is an instance of a Kubernetes Node
type nodeObj struct {
	name string
	az   string
}

var (
	podList  []podObj
	nodeList map[string]nodeObj
)

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List Kubernetes pods by AZ",
	Long: `List pods in a Kubernetes cluster by the defined availablity zone labels on the nodes
where the pods are scheduled.`,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := kube.CreateKubeClient(kubeconfig, configContext)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		nodes, label, err := kube.GetNodes(client, azLabel, cmd.Flags().Changed("label"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		nodeList := buildNodeList(nodes, label)
		podList := buildPodList(nodeList, pods)

		if len(nodeList) < 1 {
			fmt.Printf("No nodes with target AZ label (%q) found\n", label)
			os.Exit(1)
		} else if len(podList) < 1 {
			fmt.Printf("No pods were found in namespace %q\n", namespace)
			os.Exit(1)
		}

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()
		fmt.Fprintln(w, "NAME\tNAMESPACE\tNODE\tAZ\t")

		for _, pod := range podList {

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t\n", pod.name, pod.namespace, pod.node, pod.az)

		}

	},
}

func init() {
	rootCmd.AddCommand(podsCmd)
}

func buildNodeList(nodes *v1.NodeList, label string) map[string]nodeObj {

	tmpNodeList := make(map[string]nodeObj)

	for _, node := range nodes.Items {

		newNode := nodeObj{
			name: node.Name,
			az:   node.GetLabels()[label],
		}

		tmpNodeList[newNode.name] = newNode

	}

	return tmpNodeList
}

func buildPodList(nodeList map[string]nodeObj, pods *v1.PodList) []podObj {

	var tmpPodList []podObj

	for _, pod := range pods.Items {

		newPod := podObj{
			name:      pod.Name,
			namespace: pod.Namespace,
			node:      pod.Spec.NodeName,
			az:        nodeList[pod.Spec.NodeName].az,
		}

		tmpPodList = append(tmpPodList, newPod)

	}

	return tmpPodList

}
