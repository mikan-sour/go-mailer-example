package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	API_HOST              string
	API_PORT              string
	SMTP_HOST             string
	SMTP_PORT             string
	SMTP_USERNAME         string
	SMTP_PASSWORD         string
	EMAIL_FROM            string
	EMAIL_TO              string
	EMAIL_ENCRYPTION_TYPE string
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	c := &Config{
		API_HOST:              os.Getenv("API_HOST"),
		API_PORT:              os.Getenv("API_PORT"),
		SMTP_HOST:             os.Getenv("SMTP_HOST"),
		SMTP_PORT:             os.Getenv("SMTP_PORT"),
		SMTP_USERNAME:         os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD:         os.Getenv("SMTP_PASSWORD"),
		EMAIL_FROM:            os.Getenv("EMAIL_FROM"),
		EMAIL_TO:              os.Getenv("EMAIL_TO"),
		EMAIL_ENCRYPTION_TYPE: os.Getenv("EMAIL_ENCRYPTION_TYPE"),
	}

	return c, nil

}
