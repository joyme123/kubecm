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
	"github.com/spf13/cobra"
	"log"
)

var syncOpt syncOptions

type syncOptions struct {
	name string
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync refresh config from ssh",
	Long: `sync refresh config from ssh. For example:
kubecm sync k8s_cluster_1
`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := syncOptions{}
		if len(args) > 0 {
			opt.name = args[0]
		}
		sync(opt)
	},
}

func sync(opt syncOptions) {
	m, err := newManagerInterface()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	res := m.Sync(opt.name)
	for name, err := range res {
		msg := "ok"
		if err != nil {
			msg = err.Error()
		}
		log.Printf("%s: sync %s", name, msg)
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
