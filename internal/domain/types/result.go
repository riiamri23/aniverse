package types

type Result struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	AltTitles  []string `json:"altTitles"`
	Year       int      `json:"year"`
	Format     string   `json:"format"`
	Img        *string  `json:"img"`
	ProviderID string   `json:"providerId"`
}
