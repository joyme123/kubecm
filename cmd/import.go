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

type importOptions struct {
	name     string
	from string
}

var importOpt importOptions

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import config from path",
	Long:  `import config from path`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		importOpt.name = args[0]
		runImport(importOpt)
	},
}

func runImport(opt importOptions) {
	m, err := newManagerInterface()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	err = m.Import(opt.name, opt.from)
	if err != nil {
		log.Fatalf("import config error: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&importOpt.from, "from", "f", "", "import config from")
	renameCmd.MarkFlagRequired("from")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
