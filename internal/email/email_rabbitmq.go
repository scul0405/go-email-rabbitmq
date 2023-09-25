package email

import "context"

type EmailsPublisher interface {
	Publish(ctx context.Context, body []byte, contentType string) error
}

type EmailsConsumer interface {
	StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error
}