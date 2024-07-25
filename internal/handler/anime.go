package handler

import (
	"aniverse/internal/domain/types"
	"aniverse/internal/provider"
	"aniverse/internal/service"
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

type Handler struct {
	services map[string]provider.InformationProvider
}

func NewHandler(services map[string]provider.InformationProvider) *Handler {
	return &Handler{services: services}
}

func (h *Handler) GetAnimeInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	providerName := c.Query("provider", "anilist")
	var info *types.AnimeInfo

	ctx := context.Background() // Create a context

	switch providerName {
	case "anilist":
		if provider, ok := h.services["anilist"].(*service.AniList); ok {
			mediaID, err := strconv.Atoi(id)
			if err == nil {
				info, err = provider.GetAnimeInfoByID(ctx, mediaID)
				if err != nil {
					log.Printf("Error fetching anime info: %v", err)
					return c.Status(500).SendString(err.Error())
				}
			} else {
				info, err = provider.GetAnimeInfoByTitle(ctx, id)
				if err != nil {
					log.Printf("Error fetching anime info: %v", err)
					return c.Status(500).SendString(err.Error())
				}
			}
		}
	default:
		return c.Status(400).SendString("Invalid provider")
	}

	return c.JSON(info)
}

// func (h *Handler) AuthenticateKitsu(c *fiber.Ctx) error {
// 	var credentials struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	credentials.Username = os.Getenv("KITSU_USERNAME")
// 	credentials.Password = os.Getenv("KITSU_PASSWORD")

// 	if err := c.BodyParser(&credentials); err != nil {
// 		return c.Status(400).SendString("Invalid request payload")
// 	}

// 	authResponse, err := service.GetKitsuAccessToken(credentials.Username, credentials.Password)
// 	if err != nil {
// 		return c.Status(500).SendString(fmt.Sprintf("Failed to authenticate with Kitsu: %v", err))
// 	}

// 	if kitsuProvider, ok := h.services["kitsu"].(*service.Kitsu); ok {
// 		kitsuProvider.SetAccessToken(authResponse.AccessToken)
// 	}

// 	return c.JSON(authResponse)
// }
