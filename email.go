package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-gomail/gomail"
)

func EmailAlert(status TargetStatus, config Config) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Alert.FromEmail)
	msg.SetHeader("To", config.Alert.ToEmail)
	subject := "Host "
	if status.Online {
		subject += "UP: "
	} else {
		subject += "DOWN: "
	}
	subject += status.Target.Name

	statusJson, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	body := fmt.Sprintf("%s\n\n%s\n", time.Now(), statusJson)

	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	hostname := config.SMTP.Hostname
	port := config.SMTP.Port
	if config.SMTP.Hostname == "" {
		hostname = "localhost"
		port = 25
	}

	m := gomail.NewMailer(hostname, "", "", port)
	if err := m.Send(msg); err != nil {
		return fmt.Errorf("error sending alert email, err %s", err)
	}

	return nil
}
