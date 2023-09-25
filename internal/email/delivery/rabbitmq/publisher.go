package rabbitmq

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scul0405/go-email-rabbitmq/config"
	rabbitmqPkg "github.com/scul0405/go-email-rabbitmq/pkg/rabbitmq"
)

type EmailsPublisher struct {
	amqpChan *amqp.Channel
	cfg      *config.Config
}

func NewEmailsPublisher(cfg *config.Config) (*EmailsPublisher, error) {
	mqConn, err := rabbitmqPkg.NewRabbitMQConnection(cfg)
	if err != nil {
		return nil, err
	}
	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &EmailsPublisher{cfg: cfg, amqpChan: amqpChan}, nil
}

func (p *EmailsPublisher) SetupExchangeAndQueue(exchange, queueName, bindingKey, consumerTag string) error {
	err := p.amqpChan.ExchangeDeclare(
		exchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := p.amqpChan.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	err = p.amqpChan.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *EmailsPublisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		log.Printf("Close channel error: %v", err)
	}
}

func (p *EmailsPublisher) Publish(ctx context.Context, body []byte, contentType string) error {

	if err := p.amqpChan.PublishWithContext(
		ctx,
		p.cfg.RabbitMQ.Exchange,
		p.cfg.RabbitMQ.RoutingKey,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return err
	}

	return nil
}
