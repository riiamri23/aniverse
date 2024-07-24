package types

type ContentMetadata struct {
	ProviderID string      `json:"providerId"`
	Data       interface{} `json:"data"` // Can be Episode[] or Chapter[]
}
