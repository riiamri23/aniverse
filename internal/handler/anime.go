package handler

import (
	"aniverse/internal/domain/types"
	"aniverse/internal/provider"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	providers *provider.Providers
	rdb       *redis.Client
}

func NewHandler(providers *provider.Providers, rdb *redis.Client) *Handler {
	return &Handler{providers: providers, rdb: rdb}
}

func (h *Handler) GetAnimeInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	providerName := c.Query("provider", "anilist")

	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime_info:%s:%s", providerName, id)

	cached, err := h.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var info types.AnimeInfo
		if err := json.Unmarshal([]byte(cached), &info); err == nil {
			return c.JSON(info)
		}
	}

	var info *types.AnimeInfo
	if providerName == "anilist" {
		infoProvider := h.providers.AniList
		info, err = fetchInfo(infoProvider, ctx, id)
		if err != nil {
			log.Printf("Error fetching anime info: %v", err)
			return c.Status(500).SendString(err.Error())
		}
	} else {
		return c.Status(400).SendString("Invalid provider")
	}

	cacheData, _ := json.Marshal(info)
	h.rdb.Set(ctx, cacheKey, cacheData, 0)

	return c.JSON(info)
}

func (h *Handler) SearchAnime(c *fiber.Ctx) error {
	query := c.Query("query")
	providerName := c.Query("provider", "gogoanime")

	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime_search:%s:%s", providerName, query)

	cached, err := h.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var results []types.Result
		if err := json.Unmarshal([]byte(cached), &results); err == nil {
			return c.JSON(results)
		}
	}

	var results []types.Result
	if providerName == "gogoanime" {
		animeProvider := h.providers.GogoAnime
		results, err = animeProvider.Search(query)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid provider or request parameters")
	}

	cacheData, _ := json.Marshal(results)
	h.rdb.Set(ctx, cacheKey, cacheData, 0)

	return c.JSON(results)
}

func (h *Handler) GetEpisodes(c *fiber.Ctx) error {
	id := c.Params("id")
	providerName := c.Query("provider", "gogoanime")

	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime_episodes:%s:%s", providerName, id)

	cached, err := h.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var episodes []types.Episode
		if err := json.Unmarshal([]byte(cached), &episodes); err == nil {
			return c.JSON(episodes)
		}
	}

	var episodes []types.Episode
	if providerName == "gogoanime" {
		animeProvider := h.providers.GogoAnime
		episodes, err = animeProvider.Episodes(id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid provider or request parameters")
	}

	cacheData, _ := json.Marshal(episodes)
	h.rdb.Set(ctx, cacheKey, cacheData, 0)

	return c.JSON(episodes)
}

func (h *Handler) GetSource(c *fiber.Ctx) error {
	animeID := c.Params("animeID")
	epNum := c.Params("epNum")
	providerName := c.Query("provider", "gogoanime")

	if providerName != "gogoanime" {
		return c.Status(400).SendString("Invalid provider")
	}

	animeProvider := h.providers.GogoAnime

	normalizedID := strings.ToLower(strings.ReplaceAll(animeID, " ", "-"))
	episodeID := normalizedID + "-episode-" + epNum

	sources, err := animeProvider.Sources(episodeID, "sub", "gogocdn")
	if err != nil {
		log.Printf("Error fetching episode stream: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(sources)
}

func fetchInfo(provider provider.InformationProvider, ctx context.Context, identifier string) (*types.AnimeInfo, error) {
	mediaID, err := strconv.Atoi(identifier)
	if err != nil {
		return provider.GetAnimeInfoByTitle(ctx, identifier)
	}
	return provider.GetAnimeInfoByID(ctx, mediaID)
}
