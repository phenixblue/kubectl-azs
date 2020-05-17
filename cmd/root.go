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

	"kubectl-azs/pkg/kube"
	"kubectl-azs/pkg/printers"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (

	// VERSION is set during build
	VERSION string
	// Set vars for global flags
	kubeconfig    string
	configContext string
	namespace     string
	azLabel       string
	keys          []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azs",
	Short: "The \"azs\" kubectl plugin.",
	Long: `The "azs" utility is a tool to list Kubernetes objects by Availability Zone. The utility can
be used standalone or as a "kubectl" plugin`,
	Run: func(cmd *cobra.Command, args []string) {

		azs := make(map[string]string)

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

		for _, node := range nodes.Items {

			azs[node.GetLabels()[azLabel]] = node.GetLabels()[azLabel]

		}

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()
		fmt.Fprintln(w, "AZ\t")

		for az := range azs {

			keys = append(keys, az)

		}

		sort.Strings(keys)

		for _, az := range keys {

			fmt.Fprintf(w, "%v\n", azs[az])

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

func init() {
	cobra.OnInitialize()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", "", "Kubernetes configuration file")
	rootCmd.PersistentFlags().StringVar(&configContext, "context", "", "The name of the kubeconfig context to use")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "The Namespace where the proxyctl ConfigMap is located")
	rootCmd.PersistentFlags().StringVarP(&azLabel, "label", "l", "failure-domain.kubernetes.io/zone", "The target label that defines the Availability Zone on nodes")

}
