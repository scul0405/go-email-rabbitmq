package mailer

import (
	"context"

	"github.com/scul0405/go-email-rabbitmq/internal/email"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	mailDialer *gomail.Dialer
}

func NewMailer(mailDialer *gomail.Dialer) *Mailer {
	return &Mailer{mailDialer: mailDialer}
}

func (m *Mailer) Send(ctx context.Context, email *email.Email) error {
	gm := gomail.NewMessage()
	gm.SetHeader("From", email.From)
	gm.SetHeader("To", email.To...)
	gm.SetHeader("Subject", email.Subject)
	gm.SetBody(email.ContentType, email.Body)

	return m.mailDialer.DialAndSend(gm)
}
