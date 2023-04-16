package cmd

import (
	"github.com/nikhil1raghav/kindle-send/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and build info",
	Long:  `Prints the version, platform and build date for kindle-send`,
	Run: func(cmd *cobra.Command, args []string) {
		util.PrintVersion()
	},
}
