package types

type Page struct {
	URL     string            `json:"url"`
	Index   int               `json:"index"`
	Headers map[string]string `json:"headers"`
}
