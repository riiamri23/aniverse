package types

type Anime struct {
	ID                string             `json:"id"`
	Slug              string             `json:"slug"`
	CoverImage        *string            `json:"coverImage"`
	BannerImage       *string            `json:"bannerImage"`
	Trailer           *string            `json:"trailer"`
	Status            *string            `json:"status"`
	Season            string             `json:"season"`
	Title             Title              `json:"title"`
	CurrentEpisode    *int               `json:"currentEpisode"`
	Mappings          []Mapping          `json:"mappings"`
	Synonyms          []string           `json:"synonyms"`
	CountryOfOrigin   *string            `json:"countryOfOrigin"`
	Description       *string            `json:"description"`
	Duration          *int               `json:"duration"`
	Color             *string            `json:"color"`
	Year              *int               `json:"year"`
	Rating            map[string]float64 `json:"rating"`
	Popularity        map[string]float64 `json:"popularity"`
	AverageRating     *float64           `json:"averageRating,omitempty"`
	AveragePopularity *float64           `json:"averagePopularity,omitempty"`
	Type              string             `json:"type"`
	Genres            []string           `json:"genres"`
	Format            string             `json:"format"`
	Relations         []Relation         `json:"relations"`
	TotalEpisodes     *int               `json:"totalEpisodes,omitempty"`
	Episodes          Episodes           `json:"episodes"`
	Tags              []string           `json:"tags"`
	Artwork           []Artwork          `json:"artwork"`
	Characters        []Character        `json:"characters"`
}

type AnimeInfo struct {
	ID              string      `json:"id"`
	Title           Title       `json:"title"`
	Artwork         []Artwork   `json:"artwork"`
	Synonyms        []string    `json:"synonyms"`
	TotalEpisodes   *int        `json:"totalEpisodes"`
	CurrentEpisode  *int        `json:"currentEpisode"`
	BannerImage     *string     `json:"bannerImage"`
	CoverImage      *string     `json:"coverImage"`
	Color           *string     `json:"color"`
	Season          string      `json:"season"`
	Year            *int        `json:"year"`
	Status          *string     `json:"status"`
	Genres          []string    `json:"genres"`
	Description     *string     `json:"description"`
	Format          string      `json:"format"`
	Duration        *int        `json:"duration"`
	Trailer         *string     `json:"trailer"`
	CountryOfOrigin *string     `json:"countryOfOrigin"`
	Tags            []string    `json:"tags"`
	Relations       []Relation  `json:"relations"`
	Characters      []Character `json:"characters"`
	Type            string      `json:"type"`
	Rating          *float64    `json:"rating"`
	Popularity      *float64    `json:"popularity"`
}
