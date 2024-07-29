package service

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"aniverse/internal/domain/types"
// 	"aniverse/internal/helper/extractor"

// 	"github.com/PuerkitoBio/goquery"
// )

// type Anitaku struct {
// 	url      string
// 	ajaxURL  string
// 	gogoCDN  *extractor.Gogocdn
// 	streamSB *extractor.StreamSB
// }

// func NewAnitaku() *Anitaku {
// 	return &Anitaku{
// 		url:      "https://anitaku.pe",
// 		ajaxURL:  "https://ajax.gogocdn.net",
// 		gogoCDN:  extractor.NewGogocdn(nil),
// 		streamSB: extractor.NewStreamSB(nil),
// 	}
// }

// func (a *Anitaku) Search(query string) ([]types.Result, error) {
// 	results := []types.Result{}
// 	encodedQuery := url.QueryEscape(query)
// 	searchURL := fmt.Sprintf("%s/search.html?keyword=%s&page=%d", a.url, encodedQuery, page)

// 	log.Printf("Fetching search results from URL: %s", searchURL)

// 	resp, err := http.Get(searchURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}

// 	doc.Find("ul.items > li").Each(func(i int, s *goquery.Selection) {
// 		title := s.Find("p.name a").Text()
// 		id, _ := s.Find("div.img a").Attr("href")
// 		releasedText := s.Find("p.released").Text()
// 		year := 0
// 		if match := strings.TrimSpace(releasedText); match != "" {
// 			fmt.Sscanf(match, "Released: %d", &year)
// 		}
// 		img, _ := s.Find("div.img a img").Attr("src")

// 		results = append(results, types.Result{
// 			ID:         id,
// 			Title:      title,
// 			AltTitles:  []string{},
// 			Img:        &img,
// 			Format:     "TV",
// 			Year:       year,
// 			ProviderID: "anitaku",
// 		})
// 	})

// 	return results, nil
// }

// func (a *Anitaku) AdvancedSearch(keyword string, page int, season, status string) ([]types.Result, error) {
// 	if keyword == "" {
// 		return nil, fmt.Errorf("Anitaku AdvancedSearch: keyword is empty")
// 	}

// 	if page < 1 {
// 		return nil, fmt.Errorf("Anitaku AdvancedSearch: page is less than 1")
// 	}

// 	sb := strings.Builder{}
// 	sb.Grow(200)

// 	_, err := sb.WriteString(fmt.Sprintf("%s/filter.html?keyword=%s&page=%d", a.url, keyword, page))

// 	if season != "" {
// 		_, err = sb.WriteString("&season%5B%5D=" + season)
// 	}

// 	if status != "" {
// 		_, err = sb.WriteString("&status%5B%5D=" + status)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return a.Search(sb.String(), uint(page))
// }

// func (a *Anitaku) Episodes(id string) ([]types.Episode, error) {
// 	episodes := []types.Episode{}
// 	animeURL := fmt.Sprintf("%s/category/%s", a.url, id)

// 	log.Printf("Fetching episodes from URL: %s", animeURL)

// 	resp, err := http.Get(animeURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}

// 	epEnd, _ := doc.Find("#episode_page a.active").Attr("ep_end")
// 	movieID, _ := doc.Find("input#movie_id").Attr("value")
// 	defaultEp, _ := doc.Find("input#default_ep").Attr("value")
// 	alias, _ := doc.Find("input#alias_anime").Attr("value")

// 	log.Printf("epEnd: %s, movieID: %s, defaultEp: %s, alias: %s", epEnd, movieID, defaultEp, alias)

// 	ajaxURL := fmt.Sprintf("%s/ajax/load-list-episode?ep_start=%d&ep_end=%s&id=%s&default_ep=%s&alias=%s", a.ajaxURL, 0, epEnd, movieID, defaultEp, alias)

// 	log.Printf("Fetching episodes from AJAX URL: %s", ajaxURL)

// 	ajaxResp, err := http.Get(ajaxURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %w", err)
// 	}
// 	defer ajaxResp.Body.Close()

// 	if ajaxResp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("received non-200 response code: %d", ajaxResp.StatusCode)
// 	}

// 	body, _ := io.ReadAll(ajaxResp.Body)
// 	log.Printf("AJAX Response Body: %s", string(body))

// 	ajaxDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}

// 	ajaxDoc.Find("#episode_related > li").Each(func(i int, s *goquery.Selection) {
// 		epID, _ := s.Find("a").Attr("href")
// 		epNumber := strings.Replace(s.Find("div.name").Text(), "EP ", "", 1)
// 		epNum, _ := strconv.Atoi(epNumber)

// 		episodes = append(episodes, types.Episode{
// 			ID:       epID,
// 			Number:   epNum,
// 			Title:    s.Find("div.name").Text(),
// 			IsFiller: false,
// 			Img:      nil,
// 			HasDub:   strings.Contains(id, "-dub"),
// 		})
// 	})

// 	for i, j := 0, len(episodes)-1; i < j; i, j = i+1, j-1 {
// 		episodes[i], episodes[j] = episodes[j], episodes[i]
// 	}

// 	return episodes, nil
// }

// func (a *Anitaku) Servers(episodeID string) ([]types.Server, error) {
// 	if episodeID == "" {
// 		return nil, fmt.Errorf("Anitaku Servers: episode id is empty")
// 	}

// 	var servers []types.Server
// 	animeURL := fmt.Sprintf("%s/%s", a.url, episodeID)

// 	log.Printf("Fetching servers from URL: %s", animeURL)

// 	resp, err := http.Get(animeURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}

// 	doc.Find("div.anime_video_body div.anime_muti_link ul li a").Each(func(i int, s *goquery.Selection) {
// 		serverName := strings.ToLower(strings.TrimSpace(strings.Replace(s.Text(), "Choose this server", "", 1)))
// 		serverID := s.AttrOr("data-video", "")
// 		if serverID == "" {
// 			return
// 		}

// 		url, err := url.Parse(serverID)
// 		if err != nil {
// 			return
// 		}
// 		expires := url.Query().Get("expires")
// 		var expiresAt time.Time
// 		if expires != "" {
// 			seconds, _ := strconv.Atoi(expires)
// 			expiresAt = time.Unix(int64(seconds), 0)
// 		}

// 		servers = append(servers, types.Server{
// 			Name:      serverName,
// 			URL:       serverID,
// 			ExpiresAt: expiresAt,
// 		})
// 	})

// 	if len(servers) == 0 {
// 		return nil, fmt.Errorf("Anitaku Servers: no content found")
// 	}

// 	return servers, nil
// }

// func (a *Anitaku) Sources(serverID, server string) (*types.Source, error) {
// 	if !strings.HasPrefix(serverID, "http") {
// 		return nil, fmt.Errorf("Anitaku Sources: server id does not start with http")
// 	}

// 	if server == "" {
// 		return nil, fmt.Errorf("Anitaku Sources: server is empty")
// 	}

// 	var sources *types.Source
// 	var err error

// 	switch strings.ToLower(server) {
// 	case "gogo server", "gogo cdn", "gogocdn", "vidstreaming":
// 		sources, err = a.gogoCDN.Extract(serverID)
// 	case "streamsb":
// 		parsedURL, _ := url.Parse(serverID)
// 		sources, err = a.streamSB.Extract(parsedURL, false)
// 	default:
// 		return nil, fmt.Errorf("unsupported server type: %s", server)
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to extract sources: %w", err)
// 	}

// 	return sources, nil
// }
