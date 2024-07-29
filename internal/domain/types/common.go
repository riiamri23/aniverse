package types

type Title struct {
	Romaji   *string  `json:"romaji"`
	English  *string  `json:"english"`
	Native   *string  `json:"native"`
	Synonyms []string `json:"synonyms"`
}

type CoverImage struct {
	Large      string  `json:"large"`
	ExtraLarge string  `json:"extraLarge"`
	Color      *string `json:"color"`
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

type Rating struct {
	Rating            map[string]float64 `json:"rating"`
	Popularity        map[string]float64 `json:"popularity"`
	AverageRating     *float64           `json:"averageRating,omitempty"`
	AveragePopularity *float64           `json:"averagePopularity,omitempty"`
}

type Mapping struct {
	ID           string  `json:"id"`
	ProviderID   string  `json:"providerId"`
	Similarity   float64 `json:"similarity"`
	ProviderType *string `json:"providerType"`
}

type Relations struct {
	Edges []RelationEdge `json:"edges"`
}

type RelationEdge struct {
	Node         Relation `json:"node"`
	RelationType string   `json:"relationType"`
}

type Characters struct {
	Edges []CharacterEdge `json:"edges"`
}

type CharacterEdge struct {
	Node        Character    `json:"node"`
	VoiceActors []VoiceActor `json:"voiceActors"`
}
