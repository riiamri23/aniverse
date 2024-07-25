package anime

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"aniverse/internal/crawler"
	"aniverse/internal/domain/types"
	"aniverse/internal/helper/extractor"
	"aniverse/pkg/common"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Anitaku struct {
	name                string
	baseURL             string
	ajaxUrl             string
	baseCrawler         *crawler.BaseCrawler
	gogocdnExtractor    *extractor.Gogocdn
	streamwishExtractor *extractor.Streamwish
}

func NewAnitaku(baseCrawler *crawler.BaseCrawler) *Anitaku {
	if baseCrawler == nil {
		baseCrawler = crawler.NewBaseCrawler()
	}

	return &Anitaku{
		name:                "anitaku",
		baseURL:             "https://anitaku.pe",
		baseCrawler:         baseCrawler,
		ajaxUrl:             "https://ajax.gogocdn.net",
		gogocdnExtractor:    extractor.NewGogocdn(baseCrawler),
		streamwishExtractor: extractor.NewStreamwish(baseCrawler),
	}
}

func (a *Anitaku) GetName() string {
	return a.name
}

func (a *Anitaku) Search(keyword string, page uint) (*types.SearchResult, error) {
	if keyword == "" {
		return nil, fmt.Errorf("Anitaku Search: %w : keyword is empty", common.ErrInvalidArgument)
	}

	if page < 1 {
		return nil, fmt.Errorf(
			"Anitaku Search: %w : page is less than 1",
			common.ErrInvalidArgument,
		)
	}

	URL := fmt.Sprintf(
		"%s/search.html?keyword=%s&page=%d",
		a.baseURL,
		keyword,
		page,
	)

	searchResult, err := a.scrapeCardPage(URL)
	if err != nil {
		return nil, fmt.Errorf("Anitaku Search: %w", err)
	}

	return searchResult, nil
}

func (a *Anitaku) AdvanceSearch(
	keyword string,
	page int,
	season string,
	status string,
) (*types.SearchResult, error) {
	if keyword == "" {
		return nil, fmt.Errorf(
			"Anitaku AdvanceSearch: %w : keyword is empty",
			common.ErrInvalidArgument,
		)
	}

	if page < 1 {
		return nil, fmt.Errorf(
			"Anitaku AdvanceSearch: %w : page is less than 1",
			common.ErrInvalidArgument,
		)
	}

	sb := strings.Builder{}
	sb.Grow(200)

	_, err := sb.WriteString(fmt.Sprintf(
		"%s/filter.html?keyword=%s&page=%d",
		a.baseURL,
		keyword,
		page,
	))

	if season != "" {
		_, err = sb.WriteString("&season%5B%5D=" + season)
	}

	if status != "" {
		_, err = sb.WriteString("&status%5B%5D=" + status)
	}

	if err != nil {
		return nil, err
	}

	advanceSearchResult, err := a.scrapeCardPage(sb.String())
	if err != nil {
		return nil, fmt.Errorf("Anitaku AdvanceSearch: %w", err)
	}

	return advanceSearchResult, nil
}

func (a *Anitaku) Episodes(id string) (*types.Episodes, error) {
	if id == "" {
		return nil, fmt.Errorf("Anitaku Episodes: %w : id is empty", common.ErrInvalidArgument)
	}

	var episodes []types.Episode
	var errs error

	c := crawler.NewBaseCrawler()

	c.OnHTML("section.content_left div.main_body", func(h *colly.HTMLElement) {
		epEnd, err := strconv.Atoi(h.ChildAttr("#episode_page a.active", "ep_end"))
		if err != nil {
			errs = fmt.Errorf("%w : episodes info", common.ErrNoContent)
		}

		movieId := h.ChildAttr("input#movie_id", "value")
		defaultEp := h.ChildAttr("input#default_ep", "value")
		alias := h.ChildAttr("input#alias_anime", "value")

		episodes = make([]types.Episode, 0, epEnd)

		ajaxURL := fmt.Sprintf(
			"%s/ajax/load-list-episode?ep_start=%d&ep_end=%d&id=%s&default_ep=%s&alias=%s",
			a.ajaxUrl,
			0,
			epEnd,
			movieId,
			defaultEp,
			alias,
		)

		err = c.Visit(ajaxURL)
		if err != nil {
			errs = fmt.Errorf("%w : episodes info", common.ErrRequest)
			return
		}
	})

	c.OnHTML("#episode_related", func(h *colly.HTMLElement) {
		h.DOM.Find("li").Each(func(_ int, s *goquery.Selection) {
			episodeIdParts := strings.Split(s.Find("a").AttrOr("href", ""), "/")
			if len(episodeIdParts) <= 1 {
				return
			}
			episodeNo, err := strconv.Atoi(
				strings.TrimSpace(strings.Replace(s.Find("div.name").Text(), "EP", "", 1)),
			)
			if err != nil {
				return
			}

			episodes = append(episodes, types.Episode{
				ID:     episodeIdParts[1],
				Number: episodeNo,
			})
		})
	})

	URL := fmt.Sprintf("%s/category/%s", a.baseURL, id)

	err := c.Visit(URL)
	if err != nil {
		return nil, fmt.Errorf("Anitaku Episodes: %w : %s", common.ErrRequest, err.Error())
	}

	if errs != nil {
		return nil, fmt.Errorf("Anitaku Episodes: %w", errs)
	}

	if len(episodes) == 0 {
		return nil, fmt.Errorf("anitaku Episodes: %w", common.ErrNoContent)
	}

	for i, j := 0, len(episodes)-1; i < j; i, j = i+1, j-1 {
		episodes[i], episodes[j] = episodes[j], episodes[i]
	}

	return &types.Episodes{Data: []types.EpisodeData{{Episodes: episodes}}}, nil
}

func (a *Anitaku) Servers(episodeId string) (types.Servers, error) {
	if episodeId == "" {
		return nil, fmt.Errorf(
			"Anitaku Servers: %w : episode id is empty",
			common.ErrInvalidArgument,
		)
	}
	var servers types.Servers

	c := crawler.NewBaseCrawler()

	c.OnHTML("div.anime_video_body div.anime_muti_link ul", func(h *colly.HTMLElement) {
		servers = make(types.Servers, 0, h.DOM.Find("li").Length())
		isDub := strings.Contains(episodeId, "dub")

		h.DOM.Find("li a").Each(func(_ int, s *goquery.Selection) {
			serverName := strings.ToLower(
				strings.TrimSpace(strings.Replace(s.Text(), "Choose this server", "", 1)),
			)
			if serverName != "vidstreaming" && serverName != "gogo server" &&
				serverName != "streamwish" && serverName != "vidhide" {
				return
			}
			serverId := s.AttrOr("data-video", "")
			if serverId == "" {
				return
			}
			server := types.Server{
				Name: serverName,
				URL:  serverId,
			}
			url, err := url.Parse(serverId)
			if err == nil {
				expires := url.Query().Get("expires")
				if expires != "" {
					seconds, _ := strconv.Atoi(expires)
					server.ExpiresAt = time.Unix(int64(seconds), 0)
				}
			}
			if isDub {
				server.Language = "dub"
			} else {
				server.Language = "sub"
			}
			servers = append(servers, server)
		})
	})

	URL := fmt.Sprintf("%s/%s", a.baseURL, episodeId)
	err := c.Visit(URL)
	if err != nil {
		return nil, fmt.Errorf("Anitaku Servers: %w : %s", common.ErrRequest, err.Error())
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("Anitaku Servers: %w", common.ErrNoContent)
	}

	return servers, nil
}

func (a *Anitaku) Sources(serverId, server string) (*types.Source, error) {
	if !strings.HasPrefix(serverId, "http") {
		return nil, fmt.Errorf(
			"Anitaku Sources: %w : server id does not start with http",
			common.ErrInvalidArgument,
		)
	}

	if server == "" {
		return nil, fmt.Errorf("Anitaku Sources: %w : server is empty", common.ErrInvalidArgument)
	}

	var sources *types.Source

	switch strings.ToLower(server) {
	case "gogo server", "gogo cdn", "gogocdn", "vidstreaming":
		commonSources, err := a.gogocdnExtractor.Extract(serverId)
		if err != nil {
			return nil, fmt.Errorf("Anitaku Sources: %w", err)
		}
		sources = convertCommonSourcesToTypeSources(commonSources)
	case "streamwish":
		commonSources, err := a.streamwishExtractor.Extract(serverId)
		if err != nil {
			return nil, fmt.Errorf("Anitaku Sources: %w", err)
		}
		sources = convertCommonSourcesToTypeSources(commonSources)
	default:
		return nil, fmt.Errorf("Anitaku Sources: %w", common.ErrServerNotFound)
	}

	return sources, nil
}

func (a *Anitaku) scrapeCardPage(URL string) (*types.SearchResult, error) {
	cardResult := &types.SearchResult{}

	c := crawler.NewBaseCrawler()

	c.OnHTML("ul.pagination-list li", func(h *colly.HTMLElement) {
		currentPage, err := strconv.Atoi(
			h.DOM.Find("li").First().Find("a").AttrOr("data-page", ""),
		)
		if err != nil {
			currentPage = 1
		}
		cardResult.Pagination.CurrentPage = currentPage
		nextPage := h.DOM.Find("li.selected a").AttrOr("data-page", "")
		cardResult.Pagination.HasNextPage = nextPage != ""
		if !cardResult.Pagination.HasNextPage {
			cardResult.Pagination.TotalPage = currentPage
		}
	})

	c.OnHTML("ul.items", func(h *colly.HTMLElement) {
		cardResult.Data = make([]types.Search, 0, h.DOM.Find("li").Length())
		h.DOM.Find("li").Each(func(_ int, s *goquery.Selection) {
			URL := s.Find("p.name a").AttrOr("href", "")
			id := a.getId(URL)
			if id == "" {
				return
			}
			cardResult.Data = append(cardResult.Data, types.Search{
				Title: s.Find("p.name a").AttrOr("title", ""),
				Image: s.Find("a img").AttrOr("src", ""),
				Id:    id,
				Released: strings.TrimSpace(
					strings.Replace(s.Find("p.released").Text(), "Released:", "", 1),
				),
				URL: URL,
			})
		})
	})

	err := c.Visit(URL)
	if err != nil {
		return nil, fmt.Errorf("Anitaku scrapeCardPage: %w : %s", common.ErrRequest, err.Error())
	}

	if len(cardResult.Data) == 0 {
		return nil, fmt.Errorf("Anitaku scrapeCardPage: %w", common.ErrNoContent)
	}

	return cardResult, nil
}

func (a *Anitaku) getId(str string) string {
	strParts := strings.Split(str, "/")
	if len(strParts) != 3 {
		return ""
	}
	return strParts[2]
}

func convertCommonSourcesToTypeSources(commonSources *common.Sources) *types.Source {
	sources := &types.Source{
		Sources: make([]types.SourceDetail, len(commonSources.Sources)),
	}

	for i, s := range commonSources.Sources {
		sources.Sources[i] = types.SourceDetail{
			URL:     s.Url,
			Quality: s.Type,
		}
	}

	return sources
}
