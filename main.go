package main

import (
	"github.com/Pivotal-Field-Engineering/pace-builder/build"
	"github.com/Pivotal-Field-Engineering/pace-builder/serve"
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
	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Serve the PACE Workshop http://localhost:1313",
		Long:  `serve uses Hugo to serve the content.  By default Hugo uses http://localhost:1313`,
		Run: func(cmd *cobra.Command, args []string) {
			serve.ServeCmd()
		},
	}
	var rootCmd = &cobra.Command{Use: "pace"}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdServe)
	rootCmd.Execute()
}
