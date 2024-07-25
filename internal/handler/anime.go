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
	infoProvider  map[string]provider.InformationProvider
	animeProvider map[string]provider.AnimeServiceProvider
}

func NewHandler(info map[string]provider.InformationProvider, anime map[string]provider.AnimeServiceProvider) *Handler {
	return &Handler{infoProvider: info, animeProvider: anime}
}

// GetAnimeInfo fetches anime info by id or title
func (h *Handler) GetAnimeInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	providerName := c.Query("provider", "anilist")
	var info *types.AnimeInfo

	ctx := context.Background()
	provider, ok := h.infoProvider[providerName].(*service.AniList)
	if !ok {
		return c.Status(400).SendString("Invalid provider")
	}

	info, err := fetchInfo(provider, ctx, id)
	if err != nil {
		log.Printf("Error fetching anime info: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(info)
}

func (h *Handler) SearchAnime(c *fiber.Ctx) error {
	query := c.Query("q")
	providerName := c.Query("provider", "animepahe")
	var results []types.Result

	provider, ok := h.animeProvider[providerName]
	if !ok {
		return c.Status(400).SendString("Invalid provider")
	}

	results, err := provider.Search(query)
	if err != nil {
		log.Printf("Error searching anime: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(results)
}

func (h *Handler) FetchEpisodes(c *fiber.Ctx) error {
	id := c.Params("id")
	providerName := c.Query("provider", "animepahe")
	var episodes []types.Episode

	provider, ok := h.animeProvider[providerName]
	if !ok {
		return c.Status(400).SendString("Invalid provider")
	}

	episodes, err := provider.FetchEpisodes(id)
	if err != nil {
		log.Printf("Error fetching episodes: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(episodes)
}

// FetchSources fetches sources for a given episode
func (h *Handler) FetchSources(c *fiber.Ctx) error {
	id := c.Params("id")
	providerName := c.Query("provider", "animepahe")
	var source *types.Source

	provider, ok := h.animeProvider[providerName].(*service.AnimePahe)
	if !ok {
		return c.Status(400).SendString("Invalid provider")
	}

	source, err := provider.FetchSources(id, types.SubTypeSub, types.StreamingServerKwik)
	if err != nil {
		log.Printf("Error fetching sources: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(source)
}

func fetchInfo(provider *service.AniList, ctx context.Context, id string) (*types.AnimeInfo, error) {
	mediaID, err := strconv.Atoi(id)
	if err != nil {
		return provider.GetAnimeInfoByTitle(ctx, id)
	}
	return provider.GetAnimeInfoByID(ctx, mediaID)
}
