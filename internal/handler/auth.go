package handler

import (
	"aniverse/config"
	"aniverse/internal/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	tokenManager *service.TokenManager
	config       *config.Config
}

func NewAuthHandler(tokenManager *service.TokenManager, config *config.Config) *AuthHandler {
	return &AuthHandler{
		tokenManager: tokenManager,
		config:       config,
	}
}

func (h *AuthHandler) StartOAuthFlow(c *fiber.Ctx) error {
	authURL := fmt.Sprintf("https://anilist.co/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", h.config.ClientID, h.config.RedirectURI)
	return c.Redirect(authURL)
}

func (h *AuthHandler) HandleOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("No code provided")
	}

	token, err := h.tokenManager.GetAccessToken(h.config.ClientID, h.config.ClientSecret, h.config.RedirectURI, code)
	if err != nil {
		log.Printf("Failed to get access token: %v", err)
		return c.Status(500).SendString(fmt.Sprintf("Failed to get access token: %v", err))
	}

	h.tokenManager.SetToken("anilist", token.AccessToken)
	log.Printf("Access Token for AniList: %s", token.AccessToken)

	return c.JSON(token)
}
