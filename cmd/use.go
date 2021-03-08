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

	"github.com/spf13/cobra"
)

type useOptions struct {
	name string
}

var useOpt useOptions

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use set kube config file to specified one",
	Long:  `use set kube config file to specified one`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useOpt.name = args[0]
		use(useOpt)
	},
}

func use(opt useOptions) {
	m, err := newManagerInterface()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	if err := m.Use(opt.name); err != nil {
		log.Fatalf("use failed, err: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(useCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
