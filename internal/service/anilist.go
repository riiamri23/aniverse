package service

import (
	"aniverse/internal/domain/types"
	"aniverse/internal/helper"
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

const commonQueryFields = `
	id
	title {
		romaji
		english
		native
	}
	coverImage {
		extraLarge
		large
		medium
		color
	}
	description
	status
	episodes
	duration
	chapters
	volumes
	season
	seasonYear
	genres
	synonyms
	averageScore
	meanScore
	popularity
	trailer {
		id
		site
		thumbnail
	}
	bannerImage
	
`

type AniList struct {
	ID                 string
	URL                string
	API                string
	NeedsProxy         bool
	UseGoogleTranslate bool
	PreferredTitle     string
	TokenManager       *TokenManager
}

func NewAniList(tokenManager *TokenManager) *AniList {
	return &AniList{
		ID:                 "anilist",
		URL:                "https://anilist.co",
		API:                "https://graphql.anilist.co",
		NeedsProxy:         true,
		UseGoogleTranslate: false,
		PreferredTitle:     "native",
		TokenManager:       tokenManager,
	}
}

// GetAnimeInfoByID fetches information about an anime from AniList by ID
func (provider *AniList) GetAnimeInfoByID(ctx context.Context, mediaID int) (*types.AnimeInfo, error) {
	query := fmt.Sprintf(`query ($id: Int) { Media (id: $id) { %s } }`, commonQueryFields)
	return provider.fetchAnimeInfo(ctx, query, map[string]interface{}{"id": mediaID})
}

// GetAnimeInfoByTitle fetches information about an anime from AniList by title
func (provider *AniList) GetAnimeInfoByTitle(ctx context.Context, title string) (*types.AnimeInfo, error) {
	query := fmt.Sprintf(`query ($search: String) { Page(page: 1, perPage: 1) { media(search: $search, type: ANIME) { %s } } }`, commonQueryFields)
	return provider.fetchAnimeInfo(ctx, query, map[string]interface{}{"search": title})
}

// fetchAnimeInfo performs the actual data fetching from AniList API
func (provider *AniList) fetchAnimeInfo(ctx context.Context, query string, variables map[string]interface{}) (*types.AnimeInfo, error) {
	body, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := provider.new_request(ctx, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	var raw bytes.Buffer
	if _, err := raw.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Data struct {
			Media types.AnimeInfo `json:"Media"`
			Page  struct {
				Media []types.AnimeInfo `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(&raw).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	anime := &result.Data.Media
	if len(result.Data.Page.Media) > 0 {
		anime = &result.Data.Page.Media[0]
	}

	if anime == nil {
		return nil, fmt.Errorf("no anime info found")
	}

	anime.Description = helper.CleanDescription(anime.Description)
	logAnimeInfo(anime)
	return anime, nil
}

// newRequest creates a new HTTP request with the given body
func (provider *AniList) new_request(ctx context.Context, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", provider.API, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	accessToken := provider.TokenManager.GetToken("anilist")
	if accessToken == "" {
		return nil, fmt.Errorf("access token not set for AniList")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	return req, nil
}

// logAnimeInfo logs the detailed information of an anime.
func logAnimeInfo(anime *types.AnimeInfo) {
	fmt.Printf(
		"ID: %d\nTitle (Romaji): %s\nTitle (English): %s\nTitle (Native): %s\nDescription: %s\nStatus: %s\nEpisodes: %d\nDuration: %d\nSeason: %s %d\nGenres: %v\nSynonyms: %v\nAverage Score: %d\nMean Score: %d\nPopularity: %d\nBanner Image: %s\n",
		anime.ID,
		helper.GetStringValue(anime.Title.Romaji),
		helper.GetStringValue(anime.Title.English),
		helper.GetStringValue(anime.Title.Native),
		anime.Description,
		anime.Status,
		anime.Episodes,
		anime.Duration,
		anime.Season,
		anime.SeasonYear,
		anime.Genres,
		anime.Synonyms,
		anime.AverageScore,
		anime.MeanScore,
		anime.Popularity,
		anime.BannerImage,
	)
}
