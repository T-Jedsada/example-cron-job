package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jasonlvhit/gocron"
	mailgun "github.com/mailgun/mailgun-go"
)

var (
	domain        = os.Getenv("DOMAIN_MAILGUN_SANDBOX")
	apiKeyPublic  = os.Getenv("PUBLIC_KEY_MAILGUN")
	apiKeyPrivate = os.Getenv("PRIVATE_KEY_MAILGUN")
	counter       = 0
	slack         = "https://hooks.slack.com/services/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"
	contentType   = "application/json"
)

// SlackModel is model convert from file json
type SlackModel struct {
	Channel  string `json:"channel" binding:"required"`
	Username string `json:"username" binding:"required"`
	Text     string `json:"text" binding:"required"`
	IconURL  string `json:"icon_url" binding:"required"`
}

func main() {
	startJob()
}

func startJob() {
	gocron.Every(15).Minutes().Do(sendEmailWithMailGun)
	<-gocron.Start()
}

func stopJob() {
	gocron.Remove(sendEmailWithMailGun)
	gocron.Clear()
	os.Exit(1)
}

func sendEmailWithMailGun() {
	counter++
	if counter > 30 {
		stopJob()
	} else {
		gocron.Remove(sendEmailWithMailGun)
		err := sendEmail()
		if err != nil {
			fmt.Println(err)
		}
		startJob()
	}
}

func sendEmail() (err error) {
	mg := mailgun.NewMailgun(domain, apiKeyPrivate, apiKeyPublic)
	message := mg.NewMessage(
		"jt@20scoops.net",
		"test",
		fmt.Sprintf("round at : %d", counter),
		"pondthaitay@hotmail.com")
	_, _, err = mg.Send(message)
	if err != nil {
		return err
	}
	err = tickerToSlack()
	return err
}

func tickerToSlack() (err error) {
	var jsonModel SlackModel
	jsonModel.Text = fmt.Sprintf("send email round : %d", counter)
	jsonModel.Channel = "#pondthaitay_chanel"
	jsonModel.IconURL = "http://www.mailgun.com/wp-content/uploads/2017/05/mailgun.png"
	jsonModel.Username = "MailGun"
	str, err := json.Marshal(jsonModel)
	if err != nil {
		return err
	}
	body := bytes.NewBufferString(string(str))
	resp, err := http.Post(slack, contentType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return
}
