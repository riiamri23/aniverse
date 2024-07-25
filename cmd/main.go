package main

import (
	"aniverse/internal/handler"
	"aniverse/internal/provider"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := "http://localhost:3000/callback"

	var accessToken string

	app.Get("/", func(c *fiber.Ctx) error {
		authURL := fmt.Sprintf(
			"https://anilist.co/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code",
			clientID,
			redirectURI,
		)
		return c.Redirect(authURL)
	})

	app.Get("/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return c.Status(400).SendString("No code provided")
		}

		// Exchange the authorization code for an access token
		token, err := getAccessToken(clientID, clientSecret, redirectURI, code)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to get access token: %v", err))
		}

		// Store the access token
		accessToken = token.AccessToken

		// For demonstration purposes, return the access token in the response
		return c.JSON(token)
	})

	// Middleware to check if access token is set
	app.Use(func(c *fiber.Ctx) error {
		if accessToken == "" {
			return c.Status(401).SendString("Access token not set. Please authenticate first.")
		}
		return c.Next()
	})

	// Initialize provider and handler outside the callback
	prov := provider.NewProvider(accessToken)
	animeHandler := handler.NewAnimeHandler(prov.AnimeService)

	// Set up routes
	app.Get("/anime/:id", animeHandler.GetAnime)
	app.Get("/findAnime", animeHandler.FindAnime)
	app.Get("/user/:username/following", animeHandler.GetFollowingNames)
	app.Get("/user/:username/updates", animeHandler.GetUserUpdates)
	app.Get("/user/:username/progress/:mediaID", animeHandler.GetUserProgress)
	app.Post("/user/progress/update", animeHandler.UpdateUserProgress)
	log.Fatal(app.Listen(":3000"))
}

func getAccessToken(clientID, clientSecret, redirectURI, code string) (*TokenResponse, error) {
	url := "https://anilist.co/api/v2/oauth/token"
	body := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  redirectURI,
		"code":          code,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
