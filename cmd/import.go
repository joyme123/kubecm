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

	"github.com/joyme123/kubecm/pkg/loader"

	"github.com/spf13/cobra"
)

type importOptions struct {
	name       string
	from       string
	password   string
	publicKey  string
	privateKey string
	sshPort    int
}

var importOpt importOptions

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import config from path",
	Long:  `import config from path`,
	Args:  cobra.MinimumNArgs(1),
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

	var data []byte
	if loader.IsLocal(opt.from) {
		data, err = loader.LocalGet(opt.from)
		if err != nil {
			log.Fatalf("get config from local error: %v", err)
		}
	} else if loader.IsSSH(opt.from) {
		if len(opt.password) > 0 {
			data, err = loader.SSHGetWithPassword(opt.from, opt.sshPort, opt.password)
			if err != nil {
				log.Fatalf("get config from ssh with password error: %v", err)
			}
		} else {
			data, err = loader.SSHGetWithPrivateKey(opt.from, opt.sshPort, opt.privateKey)
			if err != nil {
				log.Fatalf("get config from ssh with private key error: %v", err)
			}
		}
	} else {
		log.Fatalf("unsupport path: %s", opt.from)
	}

	err = m.Import(opt.name, data)
	if err != nil {
		log.Fatalf("import config error: %v", err)
	}
}

func init() {
	privateKeyPath, err := loader.DefaultSSHKeyPath()
	if err != nil {
		log.Fatalf("get default ssh private key path error: %v", err)
	}
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&importOpt.from, "from", "f", "", "import config from")
	importCmd.MarkFlagRequired("from")
	importCmd.Flags().StringVarP(&importOpt.password, "password", "p", "", "ssh password")
	importCmd.Flags().StringVarP(&importOpt.privateKey, "privateKey", "k", privateKeyPath, "ssh private key")
	importCmd.Flags().IntVarP(&importOpt.sshPort, "port", "", 22, "ssh server port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
