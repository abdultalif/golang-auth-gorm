package workers

import (
	"encoding/json"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/services"
	"github.com/sirupsen/logrus"
)

type ForgotPasswordPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func ConsumeForgotPasswordQueue() {
	msgs, err := config.RabbitMQChannel.Consume(
		config.ForgotPasswordQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Fatalf("❌ Failed to consume forgot password queue: %v", err)
	}

	go func() {
		for d := range msgs {
			var payload ForgotPasswordPayload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				logger.Log.Errorf("❌ Invalid forgot password message: %v", err)
				continue
			}

			err = services.SendForgotPasswordEmail(payload.Email, payload.Token)
			if err != nil {
				logger.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("❌ Failed to send forgot password email")
			}
		}
	}()
}
