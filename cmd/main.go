package main

import (
	"aniverse/config"
	"aniverse/internal/handler"
	"aniverse/internal/provider"
	"aniverse/internal/service" // Import the package that contains view.WatchData
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	setting := config.LoadConfig()
	ctx := context.Background()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	tokenManager := service.NewTokenManager()
	providers := provider.NewProviders(tokenManager)
	h := handler.NewHandler(providers, rdb)
	authHandler := handler.NewAuthHandler(tokenManager, setting)

	app.Get("/", authHandler.StartOAuthFlow)
	app.Get("/callback", authHandler.HandleOAuthCallback)
	app.Use(func(c *fiber.Ctx) error {
		if tokenManager.GetToken("anilist") == "" {
			return c.Status(401).SendString("Access token not set. Please authenticate first.")
		}
		return c.Next()
	})

	app.Get("/info/:id", h.GetAnimeInfo)            // works...
	app.Get("/search", h.SearchAnime)               // works...
	app.Get("/episodes/:id", h.GetEpisodes)         // works...
	app.Get("/source/:animeID/:epNum", h.GetSource) // works ..

	// Start the server
	log.Fatal(app.Listen(setting.Port))
}
