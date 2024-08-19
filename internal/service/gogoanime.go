package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"aniverse/internal/domain/types"
	"aniverse/internal/helper/extractor"

	"github.com/PuerkitoBio/goquery"
)

type GogoAnime struct {
	url      string
	ajaxURL  string
	gogoCDN  *extractor.Gogocdn
	streamSB *extractor.StreamSB
}

func NewGogoAnime() *GogoAnime {
	return &GogoAnime{
		url:      "https://gogoanime3.co/",
		ajaxURL:  "https://ajax.gogocdn.net",
		gogoCDN:  extractor.NewGogocdn(nil),
		streamSB: extractor.NewStreamSB(nil),
	}
}

func (g *GogoAnime) Search(query string) ([]types.Result, error) {
	results := []types.Result{}
	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("%ssearch.html?keyword=%s", g.url, encodedQuery)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	doc.Find("ul.items > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find("p.name a").Text()
		id, _ := s.Find("div.img a").Attr("href")
		releasedText := s.Find("p.released").Text()
		year := 0
		if match := strings.TrimSpace(releasedText); match != "" {
			fmt.Sscanf(match, "Released: %d", &year)
		}
		img, _ := s.Find("div.img a img").Attr("src")

		results = append(results, types.Result{
			ID:         id,
			Title:      title,
			AltTitles:  []string{},
			Img:        &img,
			Format:     "TV",
			Year:       year,
			ProviderID: "gogoanime",
		})
	})

	return results, nil
}
func (g *GogoAnime) Episodes(id string) ([]types.Episode, error) {
	episodes := []types.Episode{}
	animeURL := fmt.Sprintf("%s/category/%s", g.url, id)

	resp, err := http.Get(animeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	epStart, _ := doc.Find("#episode_page > li").First().Find("a").Attr("ep_start")
	epEnd, _ := doc.Find("#episode_page > li").Last().Find("a").Attr("ep_end")
	movieID, _ := doc.Find("#movie_id").Attr("value")
	alias, _ := doc.Find("#alias_anime").Attr("value")

	ajaxURL := fmt.Sprintf("%s/ajax/load-list-episode?ep_start=%s&ep_end=%s&id=%s&default_ep=0&alias=%s", g.ajaxURL, epStart, epEnd, movieID, alias)

	ajaxResp, err := http.Get(ajaxURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer ajaxResp.Body.Close()

	if ajaxResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", ajaxResp.StatusCode)
	}

	body, _ := io.ReadAll(ajaxResp.Body)
	ajaxDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	ajaxDoc.Find("#episode_related > li").Each(func(i int, s *goquery.Selection) {
		epID, _ := s.Find("a").Attr("href")
		epNumber := strings.Replace(s.Find("div.name").Text(), "EP ", "", 1)
		epNum, _ := strconv.Atoi(epNumber)

		episodes = append(episodes, types.Episode{
			ID:       epID,
			Number:   epNum,
			Title:    s.Find("div.name").Text(),
			IsFiller: false,
			Img:      nil,
			HasDub:   strings.Contains(id, "-dub"),
		})
	})

	// Reverse the order of episodes
	for i, j := 0, len(episodes)-1; i < j; i, j = i+1, j-1 {
		episodes[i], episodes[j] = episodes[j], episodes[i]
	}

	return episodes, nil
}

func (g *GogoAnime) Sources(id, subType, server string) (*types.Source, error) {
	var serverURL string

	if strings.HasPrefix(id, "http") {
		serverURL = id
	} else {
		animeURL := fmt.Sprintf("%s%s", g.url, id)

		resp, err := http.Get(animeURL)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		switch server {
		case "GogoCDN":
			serverURL = doc.Find("#load_anime > div > div > iframe").AttrOr("src", "")
		case "VidStreaming":
			serverURL = doc.Find("div.anime_video_body > div.anime_muti_link > ul > li.vidcdn > a").AttrOr("data-video", "")
		case "StreamSB":
			serverURL = doc.Find("div.anime_video_body > div.anime_muti_link > ul > li.streamsb > a").AttrOr("data-video", "")
		default:
			serverURL = doc.Find("#load_anime > div > div > iframe").AttrOr("src", "")
		}
	}

	var sources *types.Source
	var err error

	switch {
	case strings.Contains(server, "gogocdn"):
		sources, err = g.gogoCDN.Extract(serverURL)
	case strings.Contains(server, "streamsb"):
		parsedURL, _ := url.Parse(serverURL)
		sources, err = g.streamSB.Extract(parsedURL, false)
	default:
		return nil, fmt.Errorf("unsupported server type: %s", server)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to extract sources: %w", err)
	}

	return sources, nil
}
