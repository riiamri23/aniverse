package provider

import (
	"aniverse/internal/service"
)

type Provider struct {
	AnimeService *service.AnimeService
}

func NewProvider(accessToken string) *Provider {
	return &Provider{
		AnimeService: service.NewAnimeService(accessToken),
	}
}
