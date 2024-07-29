package types

type Chapters struct {
	Latest struct {
		UpdatedAt     int64  `json:"updatedAt"`
		LatestChapter int    `json:"latestChapter"`
		LatestTitle   string `json:"latestTitle"`
	} `json:"latest"`
	Data []ChapterData `json:"data"`
}

type ChapterData struct {
	ProviderID string    `json:"providerId"`
	Chapters   []Chapter `json:"chapters"`
}

type Chapter struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Number    int      `json:"number"`
	Rating    *float64 `json:"rating"`
	UpdatedAt *int64   `json:"updatedAt,omitempty"`
}
