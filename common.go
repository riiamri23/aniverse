package common

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNoContent       = errors.New("no content")
	ErrRequest         = errors.New("request error")
	ErrJsonParse       = errors.New("JSON parsing error")
	ErrScraping        = errors.New("scraping error")
	ErrInvalidRegex    = errors.New("invalid regex")
	ErrServerNotFound  = errors.New("server not found")
)

type Flag string

const (
	FLAG_CORS_ALLOWED Flag = "CORS_ALLOWED"
)

type Source struct {
	Url    string `json:"url"`
	Type   string `json:"type"`
	IsM3U8 bool   `json:"is_m3u8"`
}

type Sources struct {
	Sources       []Source `json:"sources"`
	Thumbnail     string   `json:"thumbnail"`
	ThumbnailType string   `json:"thumbnail_type"`
	Flags         []Flag   `json:"flags"`
}
