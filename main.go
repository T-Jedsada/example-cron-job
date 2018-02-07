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
	slack         = "https://hooks.slack.com/services/" + os.Getenv("TOKEN_SLACK")
	contentType   = "application/json"
	channelSlack  = os.Getenv("CHNNALE_SLACK")
	from          = os.Getenv("FROM")
	emailTo       = os.Getenv("EMAIL_TO")
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
	gocron.Every(1).Minutes().Do(sendEmailWithMailGun)
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
		from,
		"test",
		fmt.Sprintf("round at : %d", counter),
		emailTo)
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
	jsonModel.Channel = channelSlack
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
