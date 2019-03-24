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
	"strings"

	"github.com/phenixblue/kubectl-azs/pkg/k8s"
	"github.com/phenixblue/kubectl-azs/pkg/printers"
	"github.com/spf13/cobra"
)

// podObj is an instance of a Kubernetes Pod Object
type podObj struct {
	name      string
	namespace string
	node      string
	az        string
}

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List Kubernetes pods by AZ",
	Long: `List pods in a Kubernetes cluster by the defined availablity zone labels on the nodes
where the pods are scheduled.`,
	Run: func(cmd *cobra.Command, args []string) {

		// k8sNamespace is a variable to store the namespace argument.
		var k8sNamespaceArg string

		if cmd.Flag("namespace").Value.String() != "" {

			k8sNamespaceArg = "-n" + cmd.Flag("namespace").Value.String()
		} else {

			k8sNamespaceArg = "--all-namespaces"
		}

		kubernetes := k8s.NewKubernetesCmd(true)
		out, err := kubernetes.ExecuteCommand("get", "nodes", "-l", "failure-domain.beta.kubernetes.io/zone", "-o", "custom-columns=NAME:.metadata.name,AZ:.metadata.labels.failure-domain\\.beta\\.kubernetes\\.io/zone", "--no-headers")

		if err != nil {

			fmt.Println(string(out))
			fmt.Println(err)
			os.Exit(1)

		}

		// Print stdout for debug
		//fmt.Println(string(out))

		nodeinfo := make(map[string]string)
		nodeazs := make(map[string]struct{})
		nodes := strings.Split(string(out), "\n")

		for _, node := range nodes {

			if string(node) != "" {

				kv := strings.Fields(string(node))

				k, v := kv[0], kv[1]

				// Print k,v for debug
				//fmt.Println(k, v)

				nodeinfo[k] = v
				nodeazs[v] = struct{}{}

			}

		}

		buildPods(nodeinfo, nodeazs, k8sNamespaceArg)

	},
}

func buildPods(nodeinfo map[string]string, nodeazs map[string]struct{}, k8sNamespaceCmd string) {

	podMap := make(map[string]podObj)
	var out string

	kubernetes := k8s.NewKubernetesCmd(true)

	out, err := kubernetes.ExecuteCommand("get", "pods", "-o", "custom-columns=NAME:.metadata.name,NAMESPACE:.metadata.namespace,NODE:.spec.nodeName", "--no-headers", k8sNamespaceCmd)

	if err != nil {

		fmt.Println(string(out))
		fmt.Println(err)
		os.Exit(1)

	}

	// Print stdout for debug
	//fmt.Println(string(out))

	pods := strings.Split(string(out), "\n")

	for _, pod := range pods {

		if string(pod) != "" {

			podinfo := strings.Fields(pod)

			p := podObj{
				name:      podinfo[0],
				namespace: podinfo[1],
				node:      podinfo[2],
				az:        nodeinfo[podinfo[2]],
			}

			podMap[p.name] = p

			// Print pod info for debug
			fmt.Printf("Pod: %s\n", podMap[p.name])
		}

	}

	//fmt.Println(podMap)

	w := printers.GetNewTabWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprintln(w, "NAME\tNAMESPACE\tNODE\tAZ\t")

	for az := range nodeazs {

		//fmt.Printf("# AZ: %s\n", az)

		for _, pod := range podMap {

			if pod.az == az {

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t", pod.name, pod.namespace, pod.node, pod.az)
				fmt.Fprintln(w)

			}

		}

	}

}

func init() {
	rootCmd.AddCommand(podsCmd)

	podsCmd.Flags().StringP("namespace", "n", "", "If present, the namespace scope for this CLI request")
}
