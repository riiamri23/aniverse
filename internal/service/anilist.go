package service

import (
	"aniverse/internal/domain/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
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
	query := fmt.Sprintf(`
		query ($id: Int) {
			Media (id: $id) {
				%s
			}
		}`, commonQueryFields)

	variables := map[string]interface{}{
		"id": mediaID,
	}

	return provider.fetchAnimeInfo(ctx, query, variables)
}

// GetAnimeInfoByTitle fetches information about an anime from AniList by title
func (provider *AniList) GetAnimeInfoByTitle(ctx context.Context, title string) (*types.AnimeInfo, error) {
	query := fmt.Sprintf(`
	query ($search: String) {
		Page(page: 1, perPage: 1) {
			media(search: $search, type: ANIME) {
				%s
			}
		}
	}`, commonQueryFields)

	variables := map[string]interface{}{
		"search": title,
	}

	return provider.fetchAnimeInfo(ctx, query, variables)
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

	var rawResponse bytes.Buffer
	rawResponse.ReadFrom(resp.Body)

	var result struct {
		Data struct {
			Media types.AnimeInfo `json:"Media"`
			Page  struct {
				Media []types.AnimeInfo `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(&rawResponse).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	animeInfo := &result.Data.Media
	if len(result.Data.Page.Media) > 0 {
		animeInfo = &result.Data.Page.Media[0]
	}

	if animeInfo == nil {
		return nil, fmt.Errorf("no anime info found")
	}

	animeInfo.Description = cleanDescription(animeInfo.Description)
	logAnimeInfo(animeInfo)
	return animeInfo, nil
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

// cleanDescription cleans the description text by unescaping HTML entities, removing HTML tags, and replacing newlines with spaces
func cleanDescription(description string) string {
	description = html.UnescapeString(description)
	description = removeHTMLTags(description)
	description = strings.ReplaceAll(description, "\n", " ")
	return description
}

// logAnimeInfo logs the detailed information of an anime.
func logAnimeInfo(animeInfo *types.AnimeInfo) {
	fmt.Printf(
		"ID: %d\nTitle (Romaji): %s\nTitle (English): %s\nTitle (Native): %s\nDescription: %s\nStatus: %s\nEpisodes: %d\nDuration: %d\nSeason: %s %d\nGenres: %v\nSynonyms: %v\nAverage Score: %d\nMean Score: %d\nPopularity: %d\nBanner Image: %s\n",
		animeInfo.ID,
		getStringValue(animeInfo.Title.Romaji),
		getStringValue(animeInfo.Title.English),
		getStringValue(animeInfo.Title.Native),
		animeInfo.Description,
		animeInfo.Status,
		animeInfo.Episodes,
		animeInfo.Duration,
		animeInfo.Season,
		animeInfo.SeasonYear,
		animeInfo.Genres,
		animeInfo.Synonyms,
		animeInfo.AverageScore,
		animeInfo.MeanScore,
		animeInfo.Popularity,
		animeInfo.BannerImage,
	)
}

// removeHTMLTags removes HTML tags from a string
func removeHTMLTags(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(input, "")
}

// getStringValue safely dereferences a string pointer, returning an empty string if the pointer is nil.
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
