package types

type Episodes struct {
	Latest struct {
		UpdatedAt     int64  `json:"updatedAt"`
		LatestEpisode int    `json:"latestEpisode"`
		LatestTitle   string `json:"latestTitle"`
	} `json:"latest"`
	Data []EpisodeData `json:"data"`
}

type EpisodeData struct {
	ProviderID string    `json:"providerId"`
	Episodes   []Episode `json:"episodes"`
}

type Episode struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Number      int      `json:"number"`
	IsFiller    bool     `json:"isFiller"`
	Img         *string  `json:"img"`
	HasDub      bool     `json:"hasDub"`
	Description *string  `json:"description"`
	Rating      *float64 `json:"rating"`
	UpdatedAt   *int64   `json:"updatedAt,omitempty"`
}
