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

	"kubectl-azs/pkg/k8s"

	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List Kubernetes nodes by AZ",
	Long:  `List nodes in a Kubernetes cluster by the defined availablity zone labels`,
	Run: func(cmd *cobra.Command, args []string) {

		kubernetes := k8s.NewKubernetesCmd(true)
		out, err := kubernetes.ExecuteCommand("get", "nodes", "-l", "failure-domain.beta.kubernetes.io/zone", "-o", "custom-columns=NAME:.metadata.name,AZ:.metadata.labels.failure-domain\\.beta\\.kubernetes\\.io/zone")

		if err != nil {

			fmt.Println(string(out))
			fmt.Println(err)
			os.Exit(1)

		}

		fmt.Printf(string(out))
	},
}

func init() {

	rootCmd.AddCommand(nodesCmd)

}
