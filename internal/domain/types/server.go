package types

import "time"

type Servers []Server

type Server struct {
	Server    string    `json:"server"`
	ServerId  string    `json:"serverId"`
	Language  string    `json:"language"`
	ExpiresAt time.Time `json:"expiresAt"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Type      *string   `json:"type,omitempty"`
}
