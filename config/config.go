package config

import (
	"github.com/joho/godotenv"
	"os"
	"fmt"
)

type Config struct {
	User     string
	Password string
}

func LoadConfig() (*Config, error) {

	if _, err := os.Stat(".env"); err == nil {
		// Load the .env file if it exists
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	config := &Config{
		User:     user,
		Password: password,
	}

	return config, nil
}