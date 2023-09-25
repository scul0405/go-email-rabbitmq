package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scul0405/go-email-rabbitmq/config"
)

func NewRabbitMQConnection(cfg *config.Config) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	return amqp.Dial(connAddr)
}
