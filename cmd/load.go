/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var loadOpt saveOption

type loadOption struct {
	filePath string
}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load kubeconfigs from saved file",
	Long:  `load kubeconfigs from saved file`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := newManagerInterface()
		if err != nil {
			log.Fatalf("fatal: %v", err)
		}
		err = m.Load(loadOpt.filePath)
		if err != nil {
			log.Fatalf("load failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loadCmd.Flags().StringVarP(&loadOpt.filePath, "file-path", "f", "", "load all kubeconfigs from file path")
}
