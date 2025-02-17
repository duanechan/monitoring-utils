package credentials

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EmailConfig struct {
	SMTPUser string
	SMTPPass string
	From     User
	CC       User
}

func LoadConfig() (EmailConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return EmailConfig{}, fmt.Errorf("failed to load .env file")
	}

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		return EmailConfig{}, fmt.Errorf("SMTP_USER environment variable not set")
	}

	if !IsValidEmail(smtpUser) {
		return EmailConfig{}, fmt.Errorf("invalid SMTP email")
	}

	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		return EmailConfig{}, fmt.Errorf("SMTP_PASS environment variable not set")
	}

	fromUsername := os.Getenv("SENDER_NAME")
	if fromUsername == "" {
		return EmailConfig{}, fmt.Errorf("SENDER_NAME environment variable not set")
	}

	fromEmail := os.Getenv("SENDER_EMAIL")
	if fromEmail == "" {
		return EmailConfig{}, fmt.Errorf("SENDER_EMAIL environment variable not set")
	}

	if !IsValidEmail(fromEmail) {
		return EmailConfig{}, fmt.Errorf("invalid sender email")
	}

	ccUsername := os.Getenv("CC_NAME")
	ccEmail := os.Getenv("CC_EMAIL")

	return EmailConfig{
		SMTPUser: smtpUser,
		SMTPPass: smtpPass,
		From: User{
			Name:  fromUsername,
			Email: fromEmail,
		},
		CC: User{
			Name:  ccUsername,
			Email: ccEmail,
		},
	}, nil
}
