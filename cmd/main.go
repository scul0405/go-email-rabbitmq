package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scul0405/go-email-rabbitmq/config"
	"github.com/scul0405/go-email-rabbitmq/internal/email"
	"github.com/scul0405/go-email-rabbitmq/internal/email/delivery/mailer"
	"github.com/scul0405/go-email-rabbitmq/internal/email/delivery/rabbitmq"
	"github.com/scul0405/go-email-rabbitmq/internal/email/usecase"
	mailDialerPkg "github.com/scul0405/go-email-rabbitmq/pkg/mailer"
	rabbitmqPkg "github.com/scul0405/go-email-rabbitmq/pkg/rabbitmq"
)

func main() {
	configPath := "./config/config"
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config err: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("Parse config err: %v", err)
	}

	amqpConn, err := rabbitmqPkg.NewRabbitMQConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConn.Close()

	mailDialer := mailDialerPkg.NewMailDialer(cfg)
	mailer := mailer.NewMailer(mailDialer)

	publisher, err := rabbitmq.NewEmailsPublisher(cfg)
	if err != nil {
		log.Fatalf("New publisher err: %v", err)
	}

	emailUC := usecase.NewEmailUseCase(mailer, cfg, publisher)
	emailsConsumer := rabbitmq.NewEmailsConsumer(amqpConn, emailUC)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = emailsConsumer.StartConsumer(
			cfg.RabbitMQ.WorkerPoolSize,
			cfg.RabbitMQ.Exchange,
			cfg.RabbitMQ.Queue,
			cfg.RabbitMQ.RoutingKey,
			cfg.RabbitMQ.ConsumerTag,
		)
		if err != nil {
			log.Printf("StartConsumer: %v", err)
			cancel()
		}
	}()

	// Handle server
	http.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			return
		}

		var email email.Email
		err = json.NewDecoder(r.Body).Decode(&email)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		email.EmailID = uuid.New()
		email.CreatedAt = time.Now()
		err = emailUC.PublishEmailToQueue(ctx, &email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Write([]byte("Success"))
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
