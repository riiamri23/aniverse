package types

type Title struct {
	Romaji  *string `json:"romaji"`
	English *string `json:"english"`
	Native  *string `json:"native"`
}

type Mapping struct {
	ID           string  `json:"id"`
	ProviderID   string  `json:"providerId"`
	Similarity   float64 `json:"similarity"`
	ProviderType *string `json:"providerType"`
}

type Relation struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Title        Title  `json:"title"`
	Format       string `json:"format"`
	RelationType string `json:"relationType"`
}

type Artwork struct {
	Type       string `json:"type"`
	Img        string `json:"img"`
	ProviderID string `json:"providerId"`
}

type Character struct {
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	VoiceActor VoiceActor `json:"voiceActor"`
}

type VoiceActor struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
