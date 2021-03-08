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

type renameOptions struct {
	name string
	to   string
}

var renameOpt renameOptions

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename set config name",
	Long:  `rename set config name`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		renameOpt.name = args[0]
		renameOpt.to = args[1]
		rename(renameOpt)
	},
}

func rename(opt renameOptions) {
	m, err := newManagerInterface()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	if err := m.Rename(opt.name, opt.to); err != nil {
		log.Fatalf("rename failed, err: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
