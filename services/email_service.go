package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abdultalif/golang-auth-gorm/logger"
	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(to string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	port := getIntEnv("SMTP_PORT", 587)
	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	err := d.DialAndSend(m)
	if err != nil {
		logger.Log.Errorf("❌ Failed to send email: %v", err)
		return err
	}
	return nil
}



func getIntEnv(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}



func SendForgotPasswordEmail(to string, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset Request")
	
	resetURL := os.Getenv("FRONTEND_URL") + "/reset-password?token=" + token + "&email=" + to
    htmlBody := `
        <p>Click the link below to reset your password:</p>
        <a href="` + resetURL + `" target="_blank">Reset Password</a>
    `

    m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		getIntEnv("SMTP_PORT", 587),
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	return d.DialAndSend(m)
}
