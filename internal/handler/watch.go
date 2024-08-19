package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) WatchEpisdode(c *fiber.Ctx) error {
	animeID := c.Params("animeID")
	epNum := c.Params("epNum")
	providerName := c.Query("provider", "gogoanime")

	if providerName != "gogoanime" {
		return c.Status(400).SendString("Invalid provider")
	}

	animeProvider := h.providers.GogoAnime

	normalizedID := strings.ToLower(strings.ReplaceAll(animeID, " ", "-"))
	episodeID := normalizedID + "-episode-" + epNum

	sources, err := animeProvider.Sources(episodeID, "sub", "GogoCDN")
	if err != nil {
		log.Printf("Error fetching episode stream: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(sources)
}
