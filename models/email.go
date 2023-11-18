package models

import (
	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender string = "support@zero.ua"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

func NewEmailService(config SMTPConfig) *EmailService {
	dialer := mail.NewDialer(config.Host, config.Port, config.User, config.Password)
	return &EmailService{
		dialer:        dialer,
		DefaultSender: DefaultSender,
	}
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

func (service *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	service.setFrom(email, msg)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)

	switch true {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)

	}

	return service.dialer.DialAndSend(msg)
}

func (service *EmailService) ForgotPassword(addr, url string) error {
	email := Email{
		To:        addr,
		Subject:   "Reset pw",
		Plaintext: "Just follow the fucking link: " + url,
		HTML:      `<p>Just follow the fucking : <a href="` + url + `">` + url + `</a></p>`,
	}

	return service.Send(email)
}

func (service *EmailService) setFrom(email Email, msg *mail.Message) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case service.DefaultSender != "":
		from = service.DefaultSender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}
