package types

type Source struct {
	Sources   []SourceDetail      `json:"sources"`
	Subtitles []Subtitle          `json:"subtitles"`
	Audio     []Audio             `json:"audio"`
	Intro     TimeRange           `json:"intro"`
	Outro     TimeRange           `json:"outro"`
	Headers   map[string][]string `json:"headers"`
}

type SourceDetail struct {
	URL     string `json:"url"`
	Quality string `json:"quality"`
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
