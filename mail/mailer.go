package mail

import (
	"fmt"
	"log"
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
			log.Printf("Couldn't find the file %s : %s \n", file, err)
			continue
		} else {
			tosend++
			msg.Attach(file)
		}
	}
	if tosend == 0 {
		fmt.Println("No files to send")
		return
	}

	dialer := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Sender, cfg.Password)

	fmt.Println("Sending mail")
	if err := dialer.DialAndSend(msg); err != nil {
		fmt.Println("Error sending mail ", err)
		return
	} else {
		fmt.Printf("Mailed %d files to %s", tosend, cfg.Receiver)
	}

}
