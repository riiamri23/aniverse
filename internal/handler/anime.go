package handler

import (
	"aniverse/internal/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AnimeHandler struct {
	service *service.AnimeService
}

func NewAnimeHandler(service *service.AnimeService) *AnimeHandler {
	return &AnimeHandler{service: service}
}

func (h *AnimeHandler) GetAnime(c *fiber.Ctx) error {
	id := c.Params("id")

	animeID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid anime ID")
	}

	anime, err := h.service.GetAnimeByID(animeID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if anime.ID == 0 {
		return c.Status(404).SendString("Anime not found")
	}

	return c.JSON(anime)
}

func (h *AnimeHandler) FindAnime(c *fiber.Ctx) error {
	title := c.Query("title")
	if title == "" {
		return c.Status(400).SendString("Title is required")
	}

	var firstEpisodeDate *time.Time
	if date := c.Query("firstEpisodeDate"); date != "" {
		parsedDate, err := time.Parse(time.RFC3339, date)
		if err != nil {
			return c.Status(400).SendString("Invalid date format")
		}
		firstEpisodeDate = &parsedDate
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(400).SendString("Invalid offset value")
	}

	anime, err := h.service.FindAnime(title, firstEpisodeDate, offset)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(anime)
}

func (h *AnimeHandler) GetFollowingNames(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(400).SendString("Username is required")
	}

	names, err := h.service.GetFollowingNames(username)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(names)
}

func (h *AnimeHandler) GetUserUpdates(c *fiber.Ctx) error {
	username := c.Params("username")
	mediaType := c.Query("mediaType")
	if username == "" || mediaType == "" {
		return c.Status(400).SendString("Username and mediaType are required")
	}

	var chunk *int
	if chunkParam := c.Query("chunk"); chunkParam != "" {
		chunkValue, err := strconv.Atoi(chunkParam)
		if err != nil {
			return c.Status(400).SendString("Invalid chunk value")
		}
		chunk = &chunkValue
	}

	var perChunk *int
	if perChunkParam := c.Query("perChunk"); perChunkParam != "" {
		perChunkValue, err := strconv.Atoi(perChunkParam)
		if err != nil {
			return c.Status(400).SendString("Invalid perChunk value")
		}
		perChunk = &perChunkValue
	}

	updates, err := h.service.GetUserUpdates(username, mediaType, chunk, perChunk)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(updates)
}

func (h *AnimeHandler) GetUserProgress(c *fiber.Ctx) error {
	username := c.Params("username")
	mediaIDParam := c.Params("mediaID")
	if username == "" || mediaIDParam == "" {
		return c.Status(400).SendString("Username and mediaID are required")
	}

	mediaID, err := strconv.Atoi(mediaIDParam)
	if err != nil {
		return c.Status(400).SendString("Invalid mediaID value")
	}

	progress, err := h.service.GetUserProgress(username, mediaID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"progress": progress})
}

func (h *AnimeHandler) UpdateUserProgress(c *fiber.Ctx) error {
	var payload struct {
		MediaID  int    `json:"mediaID"`
		Progress int    `json:"progress"`
		Status   string `json:"status"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).SendString("Invalid request payload")
	}

	if err := h.service.UpdateUserProgress(payload.MediaID, payload.Progress, payload.Status); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(204)
}
