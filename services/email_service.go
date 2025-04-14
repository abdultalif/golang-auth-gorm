package services

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(to string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", "Your verification code is: "+code)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		getIntEnv("SMTP_PORT", 587),
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	return d.DialAndSend(m)
}

func getIntEnv(key string, fallback int) int {
	return fallback
}
