package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scul0405/go-email-rabbitmq/internal/email"
)

type EmailsConsumer struct {
	amqpConn *amqp.Connection
	emailUC  email.UseCase
}

func NewEmailsConsumer(conn *amqp.Connection, emailUC email.UseCase) *EmailsConsumer {
	return &EmailsConsumer{
		amqpConn: conn,
		emailUC:  emailUC,
	}
}

func (c *EmailsConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *EmailsConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		err := c.emailUC.Send(ctx, delivery.Body)
		if err != nil {
			if err := delivery.Reject(false); err != nil {
				log.Printf("Err delivery.Reject: %v", err)
			}
			log.Printf("Failed to process delivery: %v", err)
		} else {
			err = delivery.Ack(false)
			if err != nil {
				log.Printf("Failed to acknowledge delivery: %v", err)
			}
		}

	}
}

func (c *EmailsConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return err
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workerPoolSize; i++ {
		c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	return chanErr
}
