package utils

import (
	"encoding/json"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func PublishMessage(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("❌ Failed to marshal data")
		return err
	}

	err = config.RabbitMQChannel.Publish(
		config.ExchangeName,
		config.RegisterRoutingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		logger.Log.Errorf("❌ Failed to publish to RabbitMQ: %v", err)
		return err
	}

	return nil
}

func PublishMessageWithRouting(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("❌ Failed to marshal data")
		return err
	}

	err = config.RabbitMQChannel.Publish(
		config.ExchangeName,
		config.ForgotPasswordRoutingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		logger.Log.Errorf("❌ Failed to publish to RabbitMQ: %v", err)
		return err
	}

	return nil
}

