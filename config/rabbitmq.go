package config

import (
	"fmt"
	"log"
	"os"

	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "auth_exchange"
	QueueName    = "email_verification"
	RoutingKey   = "email.verification"
)

var (
	RabbitMQConn    *amqp091.Connection
	RabbitMQChannel *amqp091.Channel
)

func InitRabbitMQ() {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
	os.Getenv("RABBITMQ_USERNAME"),
	os.Getenv("RABBITMQ_PASSWORD"),
	os.Getenv("RABBITMQ_HOST"),
	os.Getenv("RABBITMQ_PORT"),
	os.Getenv("RABBITMQ_VIRTUAL_HOST"),
)
conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open a channel: %v", err)
	}

	err = channel.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare exchange: %v", err)
	}

	_, err = channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type": "quorum",
		},
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare queue: %v", err)
	}

	err = channel.QueueBind(
		QueueName,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to bind queue: %v", err)
	}

	RabbitMQConn = conn
	RabbitMQChannel = channel

	logger.Log.Info("✅ Connected to RabbitMQ and declared exchange, queue, and binding")
}
