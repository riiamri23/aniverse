package types

type Manga struct {
	ID                string             `json:"id"`
	Slug              string             `json:"slug"`
	CoverImage        *string            `json:"coverImage"`
	BannerImage       *string            `json:"bannerImage"`
	Status            *string            `json:"status"`
	Title             Title              `json:"title"`
	CurrentChapter    *int               `json:"currentChapter"`
	TotalVolumes      *int               `json:"totalVolumes"`
	Color             *string            `json:"color"`
	Year              *int               `json:"year"`
	Rating            map[string]float64 `json:"rating"`
	Popularity        map[string]float64 `json:"popularity"`
	AverageRating     *float64           `json:"averageRating,omitempty"`
	AveragePopularity *float64           `json:"averagePopularity,omitempty"`
	Genres            []string           `json:"genres"`
	Type              string             `json:"type"`
	Format            string             `json:"format"`
	Relations         []Relation         `json:"relations"`
	Publisher         *string            `json:"publisher"`
	Author            *string            `json:"author"`
	TotalChapters     *int               `json:"totalChapters,omitempty"`
	Chapters          Chapters           `json:"chapters"`
	Tags              []string           `json:"tags"`
	Artwork           []Artwork          `json:"artwork"`
	Characters        []Character        `json:"characters"`
}

type MangaInfo struct {
	ID            string      `json:"id"`
	Title         Title       `json:"title"`
	Artwork       []Artwork   `json:"artwork"`
	Synonyms      []string    `json:"synonyms"`
	TotalChapters *int        `json:"totalChapters"`
	BannerImage   *string     `json:"bannerImage"`
	CoverImage    *string     `json:"coverImage"`
	Color         *string     `json:"color"`
	Year          *int        `json:"year"`
	Status        *string     `json:"status"`
	Genres        []string    `json:"genres"`
	Description   *string     `json:"description"`
	Format        string      `json:"format"`
	Publisher     *string     `json:"publisher"`
	Author        *string     `json:"author"`
	Tags          []string    `json:"tags"`
	Relations     []Relation  `json:"relations"`
	Characters    []Character `json:"characters"`
	Type          string      `json:"type"`
	Rating        *float64    `json:"rating"`
	Popularity    *float64    `json:"popularity"`
}
