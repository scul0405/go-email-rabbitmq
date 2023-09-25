package usecase

import (
	"context"
	"encoding/json"

	"github.com/scul0405/go-email-rabbitmq/config"
	"github.com/scul0405/go-email-rabbitmq/internal/email"
)

type EmailUseCase struct {
	mailer    email.Mailer
	cfg       *config.Config
	publisher email.EmailsPublisher
}

func NewEmailUseCase(mailer email.Mailer, cfg *config.Config, publisher email.EmailsPublisher) *EmailUseCase {
	return &EmailUseCase{mailer: mailer, cfg: cfg, publisher: publisher}
}

func (u *EmailUseCase) Send(ctx context.Context, deliveryBody []byte) error {
	mail := &email.Email{}
	if err := json.Unmarshal(deliveryBody, mail); err != nil {
		return err
	}

	mail.From = u.cfg.Smtp.From
	if err := u.mailer.Send(ctx, mail); err != nil {
		return err
	}

	return nil
}

func (u *EmailUseCase) PublishEmailToQueue(ctx context.Context, email *email.Email) error {
	mailBytes, err := json.Marshal(email)
	if err != nil {
		return err
	}

	return u.publisher.Publish(ctx, mailBytes, email.ContentType)
}
