package main

import (
	"aniverse/config"
	"aniverse/internal/handler"
	"aniverse/internal/provider"
	"aniverse/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	setting := config.LoadConfig()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	tokenManager := service.NewTokenManager()
	providers := provider.NewProvider(tokenManager)

	// Route to start the OAuth flow
	app.Get("/", func(c *fiber.Ctx) error {
		authURL := fmt.Sprintf("https://anilist.co/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", setting.ClientID, setting.RedirectURI)
		return c.Redirect(authURL)
	})

	// Callback route to handle OAuth response
	app.Get("/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return c.Status(400).SendString("No code provided")
		}

		token, err := getAccessToken(setting.ClientID, setting.ClientSecret, setting.RedirectURI, code)
		if err != nil {
			log.Printf("Failed to get access token: %v", err)
			return c.Status(500).SendString(fmt.Sprintf("Failed to get access token: %v", err))
		}

		tokenManager.SetToken("anilist", token.AccessToken)
		log.Printf("Access Token for AniList: %s", token.AccessToken)

		return c.JSON(token)
	})

	// Middleware to ensure access token is set
	app.Use(func(c *fiber.Ctx) error {
		if tokenManager.GetToken("anilist") == "" {
			return c.Status(401).SendString("Access token not set. Please authenticate first.")
		}
		return c.Next()
	})

	// Initialize handlers
	h := handler.NewHandler(
		map[string]provider.InformationProvider{
			"anilist": providers.AniList,
		},
		map[string]provider.AnimeServiceProvider{
			"animepahe": providers.AnimePahe,
		},
	)

	// Define routes
	app.Get("/info/:id", h.GetAnimeInfo)
	app.Get("/search", h.SearchAnime)
	app.Get("/episodes/:id", h.FetchEpisodes)
	app.Get("/sources/:id", h.FetchSources)

	// Start the server
	log.Fatal(app.Listen(setting.Port))
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
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get access token: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResponse, nil
}
