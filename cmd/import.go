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
	"github.com/joyme123/kubecm/pkg/types"

	"github.com/spf13/cobra"
)

type importOptions struct {
	name          string
	from          string
	password      string
	privateKey    string
	sshPort       int
	save          bool
	override      bool
	apiServerAddr string
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

	syncInfo := &types.Sync{
		From:       opt.from,
		SSHPort:    opt.sshPort,
		Password:   opt.password,
		PrivateKey: opt.privateKey,
	}

	data, err := loader.Load(syncInfo)
	if err != nil {
		log.Fatalf("load error: %v", err)
	}

	err = m.Import(opt.name, data, opt.override, opt.apiServerAddr)
	if err != nil {
		log.Fatalf("import config error: %v", err)
	}

	if opt.save {
		if err := m.SaveSyncInfo(opt.name, syncInfo); err != nil {
			log.Fatalf("save sync info error: %v", err)
		}
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
	importCmd.Flags().BoolVarP(&importOpt.save, "save", "", false, "save info for sync or not")
	importCmd.Flags().BoolVarP(&importOpt.override, "override", "", false, "override exist config")
	importCmd.Flags().StringVarP(&importOpt.apiServerAddr, "apiserver-addr", "", "", "apiserver address to override kubeconfig apiserver address")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
