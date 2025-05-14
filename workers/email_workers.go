package workers

import (
	"encoding/json"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/services"
	"github.com/sirupsen/logrus"
)

type VerificationPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func ConsumeVerificationQueue() {

	msgs, err := config.RabbitMQChannel.Consume(
		config.RegisterQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Fatalf("❌ Failed to consume RabbitMQ: %v", err)
	}

go func() {
	for d := range msgs {
		var payload VerificationPayload
		err := json.Unmarshal(d.Body, &payload)
		if err != nil {
			logger.Log.Errorf("❌ Invalid message format: %v", err)
			continue
		}

		err = services.SendVerificationEmail(payload.Email, payload.Code)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to send verification email")
    	}
	}
}()

}
