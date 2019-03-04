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
