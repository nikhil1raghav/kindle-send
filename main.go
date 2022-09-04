package main

import (
	"bufio"
	"flag"
	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/epubgen"
	"github.com/nikhil1raghav/kindle-send/mail"
	"log"
	"os"
)


func extractLinks(filename string) (links []string){
	links =make([]string,0)
	file, err:=os.Open(filename)
	if err!=nil{
		log.Println("Error opening link file", err)
		return
	}
	defer file.Close()
	scanner:=bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan(){
		links=append(links, scanner.Text())
	}
	return
}
func main(){
	DefaultConfig, err :=config.DefaultConfigPath()
	if err!=nil{
		log.Println(err)
		return
	}

	configPath:=flag.String("config", DefaultConfig, "Path to the configuration file (optional)")

	pageUrl:=flag.String("url", "", "URL of webpage to send")
	title:=flag.String("title", "", "Title of the epub (optional)")
	linkfile:=flag.String("linkfile", "", "Path to a text file containing multiple links separated by newline")
	filePath:=flag.String("file", "", "Mail a file to kindle, use kindle-send as a simple mailer")
	
	flag.Parse()
	passed:=0
	flag.Visit(func(f *flag.Flag){
		passed++
	})
	if passed==0{
		flag.PrintDefaults()
	}



	urls:=make([]string,0)
	if len(*pageUrl)!=0{
		urls=append(urls,*pageUrl)
	}

	if len(*linkfile)!=0{
		urls=append(urls, extractLinks(*linkfile)...)
	}

	_,err = config.Load(*configPath)
	if err!=nil{
		log.Println(err)
		return
	}
	filesToSend:=make([]string,0)
	if len(*filePath)!=0{
		filesToSend=append(filesToSend, *filePath)
	}
	if len(urls)!=0{
		book, err:=epubgen.Make(urls, *title)
		if err!=nil{
			log.Println(err)
		}else{
			filesToSend=append(filesToSend, book)
		}
	}

	if len(filesToSend)!=0{
		mail.Send(filesToSend)
	}

}
