package main

import (
	"./build"
	"github.com/spf13/cobra"
)

func main() {
	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Build the PACE Workshop",
		Long:  `build is for building a workshop based off the base PACE template, and the configuration provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			build.BuildCmd()
		},
	}
	var rootCmd = &cobra.Command{Use: "pace"}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.Execute()
}
