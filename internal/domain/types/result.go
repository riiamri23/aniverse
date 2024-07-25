package types

type Search struct {
	Title    string `json:"title"`
	Image    string `json:"image"`
	Id       string `json:"id"`
	Released string `json:"released"`
	URL      string `json:"url"`
}
type Result struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	AltTitles  []string `json:"altTitles"`
	Year       int      `json:"year"`
	Format     string   `json:"format"`
	Img        *string  `json:"img"`
	ProviderID string   `json:"providerId"`
}

type SearchResult struct {
	Pagination struct {
		CurrentPage int  `json:"currentPage"`
		TotalPage   int  `json:"totalPage"`
		HasNextPage bool `json:"hasNextPage"`
	} `json:"pagination"`
	Data []Search `json:"data"`
}
