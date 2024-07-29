package extractor

import (
	"aniverse/internal/crawler"
	"aniverse/internal/domain/types"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type StreamSB struct {
	name         string
	baseCrawler  *crawler.BaseCrawler
	reKeys       *regexp.Regexp
	reStreamData *regexp.Regexp
}

func NewStreamSB(baseCrawler *crawler.BaseCrawler) *StreamSB {
	if baseCrawler == nil {
		baseCrawler = crawler.DefaultBaseCrawler
	}

	return &StreamSB{
		name:         "streamsb",
		baseCrawler:  baseCrawler,
		reKeys:       regexp.MustCompile(`encrypted=([\da-f]+)`),
		reStreamData: regexp.MustCompile(`data-streams="(.+?)"`),
	}
}

type streamSBData struct {
	StreamData struct {
		File string `json:"file"`
	} `json:"stream_data"`
}

func (s *StreamSB) Extract(videoUrl *url.URL, isAlt bool) (*types.Source, error) {
	sources := new(types.Source)

	if videoUrl == nil {
		return nil, fmt.Errorf("StreamSB Extract: %w : video URL is not valid", ErrInvalidArgument)
	}

	id := s.extractID(videoUrl)
	if id == "" {
		return nil, fmt.Errorf("StreamSB Extract: %w : cannot extract video ID", ErrInvalidArgument)
	}

	encryptedParams, err := s.parsePage(videoUrl.String(), id)
	if err != nil {
		return nil, fmt.Errorf("StreamSB Extract: %w", err)
	}

	url := fmt.Sprintf("https://%s/sources50/%s", videoUrl.Host, encryptedParams)
	headers := map[string]string{"X-Requested-With": "XMLHttpRequest"}
	response, err := s.baseCrawler.Client.Get(url, headers)
	if err != nil {
		return nil, fmt.Errorf("StreamSB Extract: %w : %s", ErrRequest, err.Error())
	}

	var streamSBData streamSBData
	err = json.Unmarshal(response, &streamSBData)
	if err != nil {
		return nil, fmt.Errorf("StreamSB Extract: %w : stream sb data : %s", ErrJsonParse, err.Error())
	}

	if streamSBData.StreamData.File == "" {
		return nil, fmt.Errorf("StreamSB Extract: %w", ErrNoContent)
	}

	sources.Sources = append(sources.Sources, types.SourceDetail{
		URL:    streamSBData.StreamData.File,
		IsM3U8: strings.Contains(streamSBData.StreamData.File, ".m3u8"),
	})

	return sources, nil
}

func (s *StreamSB) extractID(videoUrl *url.URL) string {
	segments := strings.Split(videoUrl.Path, "/")
	if len(segments) > 1 {
		return segments[len(segments)-1]
	}
	return ""
}

func (s *StreamSB) parsePage(link string, contentID string) (string, error) {
	response, err := s.baseCrawler.Client.Get(link, nil)
	if err != nil {
		return "", fmt.Errorf("StreamSB parsePage: %w", ErrRequest)
	}

	matches := s.reKeys.FindStringSubmatch(string(response))
	if len(matches) < 2 {
		return "", fmt.Errorf("StreamSB parsePage: %w", ErrInvalidRegex)
	}
	encryptedParams := matches[1]

	streamDataMatch := s.reStreamData.FindStringSubmatch(string(response))
	if len(streamDataMatch) < 2 {
		return "", fmt.Errorf("StreamSB parsePage: %w", ErrInvalidRegex)
	}
	streamData := streamDataMatch[1]

	return fmt.Sprintf("id=%s&data=%s", encryptedParams, streamData), nil
}
