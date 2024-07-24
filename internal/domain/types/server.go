package types

type Server struct {
	Name string  `json:"name"`
	URL  string  `json:"url"`
	Type *string `json:"type,omitempty"`
}
