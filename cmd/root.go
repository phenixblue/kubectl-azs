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
	"sort"
	"strings"

	"github.com/phenixblue/kubectl-azs/pkg/k8s"
	"github.com/phenixblue/kubectl-azs/pkg/printers"
	"github.com/spf13/cobra"
)

var (

	// VERSION is set during build
	VERSION string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azs",
	Short: "The \"azs\" kubectl plugin.",
	Long: `The "azs" utility is a tool to list Kubernetes objects by Availability Zone. The utility can
be used standalone or as a "kubectl" plugin. The "kubectl" utlity needs to be installed and the "KUBECONFIG" 
environment variable needs to be set to a valid kubeconfig file.`,
	Run: func(cmd *cobra.Command, args []string) {

		kubernetes := k8s.NewKubernetesCmd(true)
		out, err := kubernetes.ExecuteCommand("get", "nodes", "-l", "failure-domain.beta.kubernetes.io/zone", "-o", "custom-columns=NAME:.metadata.name,AZ:.metadata.labels.failure-domain\\.beta\\.kubernetes\\.io/zone")

		if err != nil {

			fmt.Println(string(out))
			fmt.Println(err)
			os.Exit(1)

		}

		azs := make(map[string]struct{})
		nodes := strings.Split(string(out), "\n")

		for _, node := range nodes {

			if string(node) != "" {

				kv := strings.Split(string(node), "   ")

				v := kv[1]

				azs[v] = struct{}{}

			}

		}

		azsSort := make([]string, 0, len(azs))

		for az := range azs {
			azsSort = append(azsSort, az)
		}

		sort.Strings(azsSort)

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()
		fmt.Fprintln(w, "AZ")

		for _, az := range azsSort {

			if az != "" {
				fmt.Fprintln(w, az)
			}

		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {

	VERSION = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
