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
	goflags "flag"
	"path"

	"github.com/joyme123/kubecm/pkg/consts"
	"github.com/joyme123/kubecm/pkg/manager"
	"github.com/joyme123/kubecm/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubecm",
	Short: "kubecm manager multi kube config file",
	Long:  `kubecm manager multi kube config file`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		listCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	klog.InitFlags(nil)
	goflags.Parse()
	pflag.CommandLine.AddGoFlagSet(goflags.CommandLine)
}

func newManagerInterface() (manager.Interface, error) {
	home := util.GetHomeDir()
	configDir := path.Join(home, consts.DefaultConfigDir)
	configPath := path.Join(home, consts.DefaultConfigPath)
	kubePath := path.Join(home, consts.DefaultKubePath)
	return manager.NewInterface(configDir, configPath, kubePath)
}
