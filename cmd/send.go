package cmd

import (
	"github.com/lithammer/dedent"
	"github.com/nikhil1raghav/kindle-send/classifier"
	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/handler"
	"github.com/nikhil1raghav/kindle-send/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var (
	helpLong = `Sends the files to ereader. If a link or a file containing links is given
it will first download the webpage, convert into ebook and then send. 
Each argument is sent as a separate file.
kindle-send auto detects if argument is a link, collection of links or an ebook.`

	helpExample = dedent.Dedent(`
		# Send a single webpage
		kindle-send send "http://paulgraham.com/alien.html"

		# Send multiple webpages
		kindle-send send "http://paulgraham.com/alien.html" "http://paulgraham.com/hwh.html"

		# Send webpage, collection of webpages and an ebook
		kindle-send download "http://paulgraham.com/alien.html" links.txt "Some Book.epub"`,
	)
)

func init() {
	sendCmd.PersistentFlags().IntP("mail-timeout", "m", 120, "Mail timeout in seconds, increase it if sending lot of files")
}

var sendCmd = &cobra.Command{
	Use:     "send [LINK1] [LINK2] [FILE1] [FILE2]",
	Short:   "Send the files, links, documents to ereader",
	Long:    helpLong,
	Example: helpExample,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		_, err := config.Load(configPath)
		if err != nil {
			util.Red.Println(err)
			return
		}

		downloadRequests := classifier.Classify(args)
		downloadedRequests := handler.Queue(downloadRequests)

		timeout, err := cmd.Flags().GetInt("mail-timeout")
		if err != nil {
			timeout = 0
		}

		handler.Mail(downloadedRequests, timeout)

	},
}
