package cmd

import (
	"fmt"
	"os"

	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/util"
	"github.com/spf13/cobra"
)

func init() {
	var configPath string
	configPath, err := config.DefaultConfigPath()
	if err != nil {
		util.Red.Println("Error setting default config path: ", err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", configPath, "Path to config file")

}

var rootCmd = &cobra.Command{
	Use:   "kindle-send",
	Short: "kindle-send sends documents, webpages and books to your ereader",
	Long: `kindle-send is a CLI tool to send file (books/documents) and webpages to your ereader
It parses the webpage, optimizes it for reading on ereader, and then converts
into an ebook. Then it emails the ebook to the ereader.
Complete documentation is available at https://github.com/nikhil1raghav/kindle-send`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		_, err := config.Load(configPath)
		if err != nil {
			util.Red.Println(err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
