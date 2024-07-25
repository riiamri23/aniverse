package extractor

import (
	"aniverse/internal/crawler"
	"aniverse/pkg/common"
	"encoding/json"
	"fmt"
	"strings"
)

type Streamwish struct {
	name        string
	baseCrawler *crawler.BaseCrawler
}

func NewStreamwish(baseCrawler *crawler.BaseCrawler) *Streamwish {
	if baseCrawler == nil {
		baseCrawler = crawler.NewBaseCrawler()
	}

	return &Streamwish{
		name:        "streamwish",
		baseCrawler: baseCrawler,
	}
}

type streamwishResponse struct {
	Sources []struct {
		File string `json:"file"`
		Type string `json:"type"`
	} `json:"sources"`
}

func (s *Streamwish) Extract(link string) (*common.Sources, error) {
	sources := new(common.Sources)

	response, err := s.baseCrawler.Client.Get(link, nil)
	if err != nil {
		return nil, fmt.Errorf("Streamwish Extract: %w : %s", common.ErrRequest, err.Error())
	}

	var streamwishData streamwishResponse
	err = json.Unmarshal(response, &streamwishData)
	if err != nil {
		return nil, fmt.Errorf("Streamwish Extract: %w : streamwish data : %s", common.ErrJsonParse, err.Error())
	}

	for _, source := range streamwishData.Sources {
		if source.File == "" {
			continue
		}
		sources.Sources = append(sources.Sources, common.Source{
			Url:    source.File,
			Type:   source.Type,
			IsM3U8: strings.Contains(source.File, ".m3u8"),
		})
	}

	sources.Flags = []common.Flag{common.FLAG_CORS_ALLOWED}

	if len(sources.Sources) == 0 {
		return nil, fmt.Errorf("Streamwish Extract: %w", common.ErrNoContent)
	}

	return sources, nil
}
