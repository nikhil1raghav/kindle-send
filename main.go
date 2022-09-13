package main

import (
	"bufio"
	"flag"
	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/epubgen"
	"github.com/nikhil1raghav/kindle-send/mail"
	"github.com/nikhil1raghav/kindle-send/util"
	"os"
)

func extractLinks(filename string) (links []string) {
	links = make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		util.Red.Println("Error opening link file", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	return
}
func flagPassed(name string) bool{
	visited:=false
	flag.Visit(func(f *flag.Flag) {
		if f.Name==name{
			visited=true
		}
	})
	return visited
}
func main() {
	DefaultConfig, err := config.DefaultConfigPath()
	if err != nil {
		util.Red.Println(err)
		return
	}

	configPath := flag.String("config", DefaultConfig, "Path to the configuration file (optional)")
	pageUrl := flag.String("url", "", "URL of webpage to send")
	title := flag.String("title", "", "Title of the epub (optional)")
	linkfile := flag.String("linkfile", "", "Path to a text file containing multiple links separated by newline")
	filePath := flag.String("file", "", "Mail a file to kindle, use kindle-send as a simple mailer")
	mailTimeout :=flag.Int("mail-timeout",120, "Timeout for sending mail in Seconds" )
	_ =flag.Bool("version", false, "Print version information")
	_ =flag.Bool("dry-run", false, "Save epub locally and don't send to device")

	flag.Parse()
	passed := 0
	flag.Visit(func(f *flag.Flag) {
		passed++
	})

	if flagPassed("version"){
		util.PrintVersion()
		if passed==1{
			return
		}
	}


	if passed == 0 {
		flag.PrintDefaults()
	}


	urls := make([]string, 0)
	if len(*pageUrl) != 0 {
		urls = append(urls, *pageUrl)
	}

	if len(*linkfile) != 0 {
		urls = append(urls, extractLinks(*linkfile)...)
	}

	_, err = config.Load(*configPath)
	if err != nil {
		util.Red.Println(err)
		return
	}
	filesToSend := make([]string, 0)
	if len(*filePath) != 0 {
		filesToSend = append(filesToSend, *filePath)
	}
	if len(urls) != 0 {
		book, err := epubgen.Make(urls, *title)
		if err != nil {
			util.Red.Println(err)
		} else {
			filesToSend = append(filesToSend, book)
		}
	}

	if flagPassed("dry-run"){
		util.CyanBold.Println("Dry-run mode : Not sending files to device")
		util.Cyan.Println("Following files are saved")
		for i,file:=range filesToSend{
			util.Cyan.Printf("%d %s\n",i+1,file)
		}
		return
	}
	if len(filesToSend) != 0 {
		timeout:=config.DefaultTimeout
		if flagPassed("mail-timeout"){
			timeout=*mailTimeout
		}
		mail.Send(filesToSend, timeout)
	}

}
