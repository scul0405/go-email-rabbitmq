package email

import "context"

type UseCase interface {
	PublishEmailToQueue(ctx context.Context, email *Email) error
	Send(ctx context.Context, deliveryBody []byte) error
}