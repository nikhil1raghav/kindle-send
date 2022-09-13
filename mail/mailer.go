package mail

import (
	"github.com/nikhil1raghav/kindle-send/util"
	"os"
	"time"

	config "github.com/nikhil1raghav/kindle-send/config"

	gomail "gopkg.in/mail.v2"
)

func Send(files []string, timeout int) {
	cfg := config.GetInstance()
	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.Sender)
	msg.SetHeader("To", cfg.Receiver)

	msg.SetBody("text/plain", "")

	attachedFiles:=make([]string,0)
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			util.Red.Printf("Couldn't find the file %s : %s \n", file, err)
			continue
		} else {
			msg.Attach(file)
			attachedFiles=append(attachedFiles,file)
		}
	}
	if len(attachedFiles) == 0 {
		util.Cyan.Println("No files to send")
		return
	}

	dialer := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Sender, cfg.Password)
	dialer.Timeout=time.Duration(timeout)*time.Second
	util.CyanBold.Println("Sending mail")
	util.Cyan.Println("Mail timeout : ", dialer.Timeout.String())
	util.Cyan.Println("Following files will be sent :")
	for i,file:=range attachedFiles{
		util.Cyan.Printf("%d. %s\n",i+1,file)
	}

	if err := dialer.DialAndSend(msg); err != nil {
		util.Red.Println("Error sending mail : ", err)
		return
	} else {
		util.GreenBold.Printf("Mailed %d files to %s", len(attachedFiles), cfg.Receiver)
	}

}
