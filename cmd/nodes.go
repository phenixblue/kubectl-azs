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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List Kubernetes nodes by AZ",
	Long:  `List nodes in a Kubernetes cluster by the defined availablity zone labels`,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := kube.CreateKubeClient(kubeconfig, configContext)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: azLabel})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(nodes.Items) < 1 {
			fmt.Printf("No nodes with target AZ label (%q) found\n", azLabel)
			os.Exit(1)
		}

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()
		fmt.Fprintln(w, "NODE NAME\tAZ\t")

		for _, node := range nodes.Items {

			fmt.Fprintf(w, "%v\t%v\t\n", node.GetName(), node.GetLabels()[azLabel])

		}

	},
}

func init() {

	rootCmd.AddCommand(nodesCmd)

}
