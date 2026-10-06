package config

import (
	"log"
	"os"
)

type ResendConfig struct {
	APIKey    string
	EmailFrom string
	ReplyTo   string
}

func LoadResendConfigFromEnv() *ResendConfig {
	LoadEnv()

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Fatal("RESEND_API_KEY environment variable is not set")
	}

	emailFrom := os.Getenv("EMAIL_FROM")
	if apiKey == "" {
		log.Fatal("EMAIL_FROM environment variable is not set")
	}

	replyTo := os.Getenv("REPLY_TO")
	if apiKey == "" {
		log.Fatal("REPLY_TO environment variable is not set")
	}

	return &ResendConfig{
		APIKey:    apiKey,
		EmailFrom: emailFrom,
		ReplyTo:   replyTo,
	}
}
