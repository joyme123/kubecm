/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type listOption struct {
	Current bool
}

var listOpt listOption

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all configurations",
	Long:  `list all configurations`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func list() {
	m, err := newManagerInterface()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}

	startTime := time.Now()

	c, err := m.List()
	if err != nil {
		log.Fatalf("list configs error: %v", err)
	}
	c.Print(listOpt.Current)
	cost := time.Since(startTime)
	klog.V(4).Infof("list cost: %d ms", cost.Milliseconds())
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listOpt.Current, "current", "c", false, "only show current kube config name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
