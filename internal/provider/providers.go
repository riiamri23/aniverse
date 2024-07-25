package provider

import (
	"aniverse/internal/service"
)

// Provider is a collection of service providers
type Provider struct {
	TokenManager *service.TokenManager
	AniList      InformationProvider
	AnimePahe    AnimeServiceProvider
	// Kitsu        InfoProvider
	// Inject other service providers here
}

// NewProvider creates a new provider
func NewProvider(tokenManager *service.TokenManager) *Provider {
	return &Provider{
		TokenManager: tokenManager,
		AniList:      service.NewAniList(tokenManager),
		AnimePahe:    service.NewAnimePahe(),
		// Kitsu:        service.NewKitsu(tokenManager),
		// Instantiate other providers here
	}
}
