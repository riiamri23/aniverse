package types

type Anime struct {
	ID                string             `json:"id"`
	Slug              string             `json:"slug"`
	CoverImage        *CoverImage        `json:"coverImage"`
	BannerImage       string             `json:"bannerImage"`
	Trailer           string             `json:"trailer"`
	Status            string             `json:"status"`
	Season            string             `json:"season"`
	Title             Title              `json:"title"`
	CurrentEpisode    int                `json:"currentEpisode"`
	Mappings          []Mapping          `json:"mappings"`
	Synonyms          []string           `json:"synonyms"`
	CountryOfOrigin   string             `json:"countryOfOrigin"`
	Description       string             `json:"description"`
	Duration          int                `json:"duration"`
	Color             string             `json:"color"`
	Year              int                `json:"year"`
	Rating            map[string]float64 `json:"rating"`
	Popularity        map[string]float64 `json:"popularity"`
	AverageRating     float64            `json:"averageRating,omitempty"`
	AveragePopularity float64            `json:"averagePopularity,omitempty"`
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
	ID           int        `json:"id"`
	Title        Title      `json:"title"`
	CoverImage   CoverImage `json:"coverImage"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Episodes     int        `json:"episodes"`
	Duration     int        `json:"duration"`
	Chapters     *int       `json:"chapters,omitempty"`
	Volumes      *int       `json:"volumes,omitempty"`
	Season       string     `json:"season"`
	SeasonYear   int        `json:"seasonYear"`
	Genres       []string   `json:"genres"`
	Synonyms     []string   `json:"synonyms"`
	AverageScore int        `json:"averageScore"`
	MeanScore    int        `json:"meanScore"`
	Popularity   int        `json:"popularity"`
	Trailer      *Trailer   `json:"trailer,omitempty"`
	BannerImage  string     `json:"bannerImage"`
}
