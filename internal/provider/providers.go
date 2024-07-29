package provider

import (
	"aniverse/internal/service" // Import the package that contains NewGoGoAnime
)

// Provider is a collection of service providers
type Providers struct {
	AniList      InformationProvider
	GogoAnime    AnimeServiceProvider
	Anitaku      AnimeServiceProvider
	TokenManager *service.TokenManager
}

// NewProvider creates a new provider
func NewProviders(tokenManager *service.TokenManager) *Providers {
	return &Providers{
		TokenManager: tokenManager,
		AniList:      service.NewAniList(tokenManager),
		GogoAnime:    service.NewGogoAnime(),
		// Anitaku:      service.NewAnitaku(),
	}
}
