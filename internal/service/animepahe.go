package service

import (
	"aniverse/internal/crawler"
	"aniverse/internal/domain/types"
	"aniverse/internal/helper/extractor"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AnimePahe struct {
	rateLimit          int
	id                 string
	url                string
	needsProxy         bool
	useGoogleTranslate bool
	formats            []types.Format
	headers            map[string]string
	baseCrawler        *crawler.BaseCrawler
}

func NewAnimePahe() *AnimePahe {
	return &AnimePahe{
		rateLimit:          250,
		id:                 "animepahe",
		url:                "https://animepahe.ru",
		needsProxy:         true,
		useGoogleTranslate: false,
		formats:            []types.Format{types.FormatMovie, types.FormatONA, types.FormatOVA, types.FormatSpecial, types.FormatTV, types.FormatTVShort},
		headers:            map[string]string{"Referer": "https://kwik.si"},
		baseCrawler:        crawler.NewBaseCrawler(),
	}
}

func (a *AnimePahe) Search(query string) ([]types.Result, error) {
	searchURL := fmt.Sprintf("%s/api?m=search&q=%s", a.url, url.QueryEscape(query))
	headers := map[string]string{"Cookie": "__ddg1_=;__ddg2_=;"}
	client := a.baseCrawler.Client
	resp, err := client.Get(searchURL, headers)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Year    int    `json:"year"`
			Poster  string `json:"poster"`
			Type    string `json:"type"`
			Session string `json:"session"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	results := []types.Result{}
	for _, item := range data.Data {
		formatString := strings.ToUpper(item.Type)
		format := types.FormatUnknown
		for _, f := range a.formats {
			if string(f) == formatString {
				format = f
				break
			}
		}

		results = append(results, types.Result{
			ID:         fmt.Sprintf("%d-%s", item.ID, item.Session),
			Title:      item.Title,
			Year:       item.Year,
			Img:        &item.Poster,
			Format:     string(format),
			ProviderID: a.id,
		})
	}

	return results, nil
}

func (a *AnimePahe) FetchEpisodes(id string) ([]types.Episode, error) {
	apiURL := fmt.Sprintf("https://api.anify.tv/episodes/%s", id)
	client := a.baseCrawler.Client
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", "your_token_here"), // Replace with actual token
	}
	resp, err := client.Get(apiURL, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching episodes: %w", err)
	}

	var externalResponse []types.EpisodeData
	if err := json.Unmarshal(resp, &externalResponse); err != nil {
		return nil, fmt.Errorf("error decoding episodes response: %w", err)
	}

	if len(externalResponse) == 0 {
		return nil, fmt.Errorf("no episodes found for anime id: %s", id)
	}

	episodes := externalResponse[0].Episodes
	return episodes, nil
}

func (a *AnimePahe) FetchSources(id string, subType types.SubType, server types.StreamingServer) (*types.Source, error) {
	splitID := strings.LastIndex(id, "-")
	if splitID == -1 {
		return nil, fmt.Errorf("invalid ID format")
	}

	animeID := id[splitID+1:]
	episodeID := id[:splitID]

	animeURL := fmt.Sprintf("%s/%s", a.url, animeID)
	client := a.baseCrawler.Client
	resp, err := client.Get(animeURL, nil)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, err
	}

	tempID := doc.Find("head > meta[property='og:url']").AttrOr("content", "")
	tempID = tempID[strings.LastIndex(tempID, "/")+1:]

	apiURL := fmt.Sprintf("%s/api?m=release&id=%s&sort=episode_asc&page=1", a.url, tempID)
	resp, err = client.Get(apiURL, map[string]string{"Cookie": "__ddg1_=;__ddg2_=;"})
	if err != nil {
		return nil, err
	}

	var data struct {
		LastPage int `json:"last_page"`
		Data     []struct {
			ID      int    `json:"id"`
			Session string `json:"session"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	episodeSession := ""
	for _, item := range data.Data {
		if fmt.Sprintf("%d", item.ID) == episodeID {
			episodeSession = item.Session
			break
		}
	}

	if episodeSession == "" {
		for i := 1; i < data.LastPage; i++ {
			pageURL := fmt.Sprintf("%s/api?m=release&id=%s&sort=episode_asc&page=%d", a.url, tempID, i+1)
			resp, err = client.Get(pageURL, map[string]string{"Cookie": "__ddg1_=;__ddg2_=;"})
			if err != nil {
				return nil, err
			}

			var pageData struct {
				Data []struct {
					ID      int    `json:"id"`
					Session string `json:"session"`
				} `json:"data"`
			}
			if err := json.Unmarshal(resp, &pageData); err != nil {
				return nil, err
			}

			for _, item := range pageData.Data {
				if fmt.Sprintf("%d", item.ID) == episodeID {
					episodeSession = item.Session
					break
				}
			}

			if episodeSession != "" {
				break
			}
		}
	}

	if episodeSession == "" {
		return nil, fmt.Errorf("episode session not found")
	}

	watchURL := fmt.Sprintf("%s/play/%s/%s", a.url, animeID, episodeSession)
	resp, err = client.Get(watchURL, map[string]string{"Cookie": "__ddg1_=;__ddg2_=;"})
	if err != nil {
		return nil, err
	}
	htmlData := string(resp)

	regex := regexp.MustCompile(`https:\/\/kwik\.si\/e\/\w+`)
	matches := regex.FindAllString(htmlData, -1)
	if matches == nil {
		return nil, fmt.Errorf("no matches found for kwik.si URL")
	}

	extractor := extractor.NewGogocdn(a.baseCrawler)
	commonSource, err := extractor.Extract(matches[0])
	if err != nil {
		return nil, err
	}

	source := &types.Source{
		Sources:   []types.SourceDetail{},
		Subtitles: []types.Subtitle{},
		Audio:     []types.Audio{},
		Intro:     types.TimeRange{Start: 0, End: 0},
		Outro:     types.TimeRange{Start: 0, End: 0},
		Headers:   a.headers,
	}

	for _, s := range commonSource.Sources {
		source.Sources = append(source.Sources, types.SourceDetail{
			URL:     s.Url,
			Quality: s.Type,
		})
	}

	return source, nil
}
