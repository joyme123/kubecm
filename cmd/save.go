/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var saveOpt saveOption

type saveOption struct {
	filePath string
}

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "save all kubeconfigs into a single file",
	Long:  `you can use save command to export all kubeconfigs into a single file. And also you can use load command to import all of these kubeconfigs`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := newManagerInterface()
		if err != nil {
			log.Fatalf("fatal: %v", err)
		}
		err = m.Save(saveOpt.filePath)
		if err != nil {
			log.Fatalf("save failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	saveCmd.Flags().StringVarP(&saveOpt.filePath, "file-path", "f", "", "saved file path for all kubeconfigs")
}
