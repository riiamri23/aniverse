package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port         string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	clientID := os.Getenv("ANILIST_CLIENT_ID")
	clientSecret := os.Getenv("ANILIST_CLIENT_SECRET")
	redirectURI := os.Getenv("ANILIST_REDIRECT_URI")

	log.Printf("Loaded environment variables: PORT=%s, ANILIST_CLIENT_ID=%s, ANILIST_CLIENT_SECRET=%s, ANILIST_REDIRECT_URI=%s",
		port, clientID, clientSecret, redirectURI)

	return &Config{
		Port:         port,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}
