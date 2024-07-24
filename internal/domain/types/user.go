package types

type Key struct {
	ID           string `json:"id"`
	Key          string `json:"key"`
	RequestCount int    `json:"requestCount"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}
