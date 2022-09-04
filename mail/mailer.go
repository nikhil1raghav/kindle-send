package mail

import (
	"github.com/nikhil1raghav/kindle-send/util"
	"os"

	config "github.com/nikhil1raghav/kindle-send/config"

	gomail "gopkg.in/mail.v2"
)

func Send(files []string) {
	cfg := config.GetInstance()
	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.Sender)
	msg.SetHeader("To", cfg.Receiver)

	msg.SetBody("text/plain", "")

	tosend := 0
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			util.Red.Printf("Couldn't find the file %s : %s \n", file, err)
			continue
		} else {
			tosend++
			msg.Attach(file)
		}
	}
	if tosend == 0 {
		util.Cyan.Println("No files to send")
		return
	}

	dialer := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Sender, cfg.Password)

	util.CyanBold.Println("Sending mail")
	if err := dialer.DialAndSend(msg); err != nil {
		util.Red.Println("Error sending mail : ", err)
		return
	} else {
		util.GreenBold.Printf("Mailed %d files to %s", tosend, cfg.Receiver)
	}


}
