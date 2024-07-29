package main

import (
	"aniverse/config"
	"aniverse/internal/handler"
	"aniverse/internal/provider"
	"aniverse/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	setting := config.LoadConfig()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	tokenManager := service.NewTokenManager()
	providers := provider.NewProviders(tokenManager)
	h := handler.NewHandler(providers)
	authHandler := handler.NewAuthHandler(tokenManager, setting)

	app.Get("/", authHandler.StartOAuthFlow)
	app.Get("/callback", authHandler.HandleOAuthCallback)
	app.Use(func(c *fiber.Ctx) error {
		if tokenManager.GetToken("anilist") == "" {
			return c.Status(401).SendString("Access token not set. Please authenticate first.")
		}
		return c.Next()
	})

	app.Get("/info/:id", h.GetAnimeInfo)
	app.Get("/search", h.SearchAnime)
	app.Get("/episodes", h.GetEpisodes)
	app.Get("/watch/:animeID/:epNum", h.GetEpisodeStream)
	app.Get("/sources", h.GetSources)

	// Start the server
	log.Fatal(app.Listen(setting.Port))
}
