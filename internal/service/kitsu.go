package service

// import (
// 	"aniverse/internal/domain/types"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// )

// type Kitsu struct {
// 	ID                 string
// 	URL                string
// 	API                string
// 	NeedsProxy         bool
// 	UseGoogleTranslate bool
// 	PreferredTitle     string
// 	AccessToken        string
// }

// type KitsuResponse struct {
// 	Data struct {
// 		Attributes struct {
// 			Titles struct {
// 				EN   string `json:"en"`
// 				ENJP string `json:"en_jp"`
// 				JAJP string `json:"ja_jp"`
// 			} `json:"titles"`
// 			Description   string `json:"description"`
// 			Subtype       string `json:"subtype"`
// 			Status        string `json:"status"`
// 			ShowType      string `json:"showType"`
// 			Synopsis      string `json:"synopsis"`
// 			EpisodeLength int    `json:"episodeLength"`
// 			PosterImage   struct {
// 				Original string `json:"original"`
// 			} `json:"posterImage"`
// 			CoverImage struct {
// 				Original string `json:"original"`
// 			} `json:"coverImage"`
// 			AverageRating string `json:"averageRating"`
// 			EpisodeCount  int    `json:"episodeCount"`
// 		} `json:"attributes"`
// 	} `json:"data"`
// }

// func NewKitsu() *Kitsu {
// 	return &Kitsu{
// 		ID:             "kitsu",
// 		URL:            "https://kitsu.io",
// 		API:            "https://kitsu.io/api/edge",
// 		NeedsProxy:     true,
// 		PreferredTitle: "native",
// 	}
// }

// func (k *Kitsu) SetAccessToken(token string) {
// 	k.AccessToken = token
// }

// func (k *Kitsu) GetAnimeInfo(mediaID int) (*types.AnimeInfo, error) {
// 	url := fmt.Sprintf("%s/anime/%d", k.API, mediaID)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Accept", "application/vnd.api+json")
// 	req.Header.Set("Content-Type", "application/vnd.api+json")
// 	req.Header.Set("Authorization", "Bearer "+k.AccessToken)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var kitsuResponse KitsuResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&kitsuResponse); err != nil {
// 		return nil, err
// 	}

// 	attributes := kitsuResponse.Data.Attributes
// 	genres, err := k.getGenres(mediaID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	artwork := []types.Artwork{}
// 	if attributes.CoverImage.Original != "" {
// 		artwork = append(artwork, types.Artwork{
// 			Type:       "banner",
// 			Img:        attributes.CoverImage.Original,
// 			ProviderID: k.ID,
// 		})
// 	}

// 	if attributes.PosterImage.Original != "" {
// 		artwork = append(artwork, types.Artwork{
// 			Type:       "poster",
// 			Img:        attributes.PosterImage.Original,
// 			ProviderID: k.ID,
// 		})
// 	}

// 	return &types.AnimeInfo{
// 		ID: fmt.Sprintf("%d", mediaID),
// 		Title: &types.Title{
// 			Romaji:  &attributes.Titles.ENJP,
// 			English: &attributes.Titles.EN,
// 			Native:  &attributes.Titles.JAJP},
// 		CurrentEpisode:  nil,
// 		Trailer:         nil,
// 		Duration:        &attributes.EpisodeLength,
// 		Color:           nil,
// 		BannerImage:     &attributes.CoverImage.Original,
// 		CoverImage:      &attributes.PosterImage.Original,
// 		Status:          &attributes.Status,
// 		Format:          string(types.FormatUnknown),
// 		Season:          string(types.SeasonUnknown),
// 		Synonyms:        []string{},
// 		Description:     &attributes.Synopsis,
// 		Year:            nil,
// 		TotalEpisodes:   &attributes.EpisodeCount,
// 		Genres:          genres,
// 		Rating:          parseRating(attributes.AverageRating),
// 		Popularity:      0,
// 		CountryOfOrigin: nil,
// 		Tags:            []string{},
// 		Relations:       []types.Relation{},
// 		Artwork:         artwork,
// 		Characters:      []types.Character{},
// 	}, nil
// }

// func (k *Kitsu) getGenres(mediaID int) ([]string, error) {
// 	url := fmt.Sprintf("%s/anime/%d/genres", k.API, mediaID)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+k.AccessToken)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var response struct {
// 		Data []struct {
// 			Attributes struct {
// 				Name string `json:"name"`
// 			} `json:"attributes"`
// 		} `json:"data"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return nil, err
// 	}

// 	genres := make([]string, len(response.Data))
// 	for i, genre := range response.Data {
// 		genres[i] = genre.Attributes.Name
// 	}

// 	return genres, nil
// }

// func parseRating(rating string) float64 {
// 	if rating == "" {
// 		return 0
// 	}
// 	r, err := strconv.ParseFloat(rating, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return r / 10
// }
