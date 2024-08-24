package handler

import (
	"aniverse/internal/domain/types"
	"aniverse/view"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) WatchEpisode(c *fiber.Ctx) error {
	c.Type("html")

	animeID := c.Params("animeID")
	epNum := c.Params("epNum")
	providerName := c.Query("provider", "gogoanime")

	if providerName != "gogoanime" {
		return c.Status(400).SendString("Invalid provider")
	}

	provider := h.providers.GogoAnime
	normalizedID := strings.ToLower(strings.ReplaceAll(animeID, " ", "-"))
	episodeID := normalizedID + "-episode-" + epNum

	// Fetch episode stream data
	data, err := provider.Sources(episodeID, "sub", "gogocdn")
	if err != nil {
		log.Printf("Error fetching episode stream: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	source := types.Source{
		URL:           data.URL,
		Type:          data.Type,
		IsM3U8:        true,
		Thumbnail:     data.Thumbnail,
		ThumbnailType: data.ThumbnailType,
		Flags:         nil,
		Subtitles:     nil,
		Audio:         nil,
	}

	return view.Watch(&source).Render(c.Context(), c)
}
