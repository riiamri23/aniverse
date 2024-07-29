package handler

import (
	"aniverse/internal/domain/types"
	"aniverse/internal/provider"
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	providers *provider.Providers
}

func NewHandler(providers *provider.Providers) *Handler {
	return &Handler{providers: providers}
}

func (h *Handler) GetAnimeInfo(c *fiber.Ctx) error {
	identifier := c.Params("id")
	providerName := c.Query("provider", "anilist")

	ctx := context.Background()

	var info *types.AnimeInfo
	var err error

	if providerName == "anilist" {
		infoProvider := h.providers.AniList
		info, err = fetchInfo(infoProvider, ctx, identifier)
	} else {
		return c.Status(400).SendString("Invalid provider")
	}

	if err != nil {
		log.Printf("Error fetching anime info: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(info)
}

func (h *Handler) SearchAnime(c *fiber.Ctx) error {
	providerName := c.Query("provider", "gogoanime")
	query := c.Query("query")

	var results []types.Result
	var err error

	if providerName == "gogoanime" {
		animeProvider := h.providers.GogoAnime
		results, err = animeProvider.Search(query)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid provider or request parameters")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(results)
}

func (h *Handler) GetEpisodes(c *fiber.Ctx) error {
	providerName := c.Query("provider", "gogoanime")
	id := c.Query("id")

	var episodes []types.Episode
	var err error

	if providerName == "gogoanime" {
		animeProvider := h.providers.GogoAnime
		episodes, err = animeProvider.Episodes(id)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid provider or request parameters")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(episodes)
}

func (h *Handler) GetSources(c *fiber.Ctx) error {
	providerName := c.Query("provider", "gogoanime")
	id := c.Query("id")
	subType := c.Query("subType")
	server := c.Query("server")

	var sources *types.Source
	var err error

	if providerName == "gogoanime" {
		animeProvider := h.providers.GogoAnime
		sources, err = animeProvider.Sources(id, subType, server)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid provider or request parameters")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(sources)
}

// New handler function for the /:animeID/:epNum route
func (h *Handler) GetEpisodeStream(c *fiber.Ctx) error {
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

func fetchInfo(provider provider.InformationProvider, ctx context.Context, identifier string) (*types.AnimeInfo, error) {
	mediaID, err := strconv.Atoi(identifier)
	if err != nil {
		// Identifier is a title
		return provider.GetAnimeInfoByTitle(ctx, identifier)
	}
	// Identifier is an ID
	return provider.GetAnimeInfoByID(ctx, mediaID)
}
