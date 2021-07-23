package config

import (
	"os"
)

const (
	TelegramApiUrl = "https://api.telegram.org"
	Domain         = "https://a2ed86a981f7.ngrok.io"
)

type TelegramConfig struct {
	Token string
}

type Config struct {
	Telegram TelegramConfig
}

// NewConfig returns a new Config struct
func NewConfig() *Config {
	return &Config{
		Telegram: TelegramConfig{
			Token: getEnv("TOKEN", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
