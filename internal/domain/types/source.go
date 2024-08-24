package types

type Source struct {
	URL           string              `json:"url"`
	Type          string              `json:"type"`
	IsM3U8        bool                `json:"is_m3u8"`
	Thumbnail     string              `json:"thumbnail"`
	ThumbnailType string              `json:"thumbnail_type"`
	Flags         []Flag              `json:"flags"`
	Subtitles     []Subtitle          `json:"subtitles"`
	Audio         []Audio             `json:"audio"`
	Intro         TimeRange           `json:"intro"`
	Outro         TimeRange           `json:"outro"`
	Headers       map[string][]string `json:"headers"`
}

type Subtitle struct {
	URL   string `json:"url"`
	Lang  string `json:"lang"`
	Label string `json:"label,omitempty"`
}

type Audio struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	Language string `json:"language"`
}
