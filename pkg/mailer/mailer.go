package mailer

import (
	"github.com/scul0405/go-email-rabbitmq/config"
	"gopkg.in/gomail.v2"
)

func NewMailDialer(cfg *config.Config) *gomail.Dialer {
	return gomail.NewDialer(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.User, cfg.Smtp.Password)
}
