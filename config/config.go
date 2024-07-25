package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		ClientID:     os.Getenv("ANILIST_CLIENT_ID"),
		ClientSecret: os.Getenv("ANILIST_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
	}
}
