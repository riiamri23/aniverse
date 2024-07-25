package provider

import (
	"aniverse/internal/domain/types"
	"context"
)

// InformationProvider is an interface for fetching anime information
type InformationProvider interface {
	GetAnimeInfoByID(ctx context.Context, mediaID int) (*types.AnimeInfo, error)
	GetAnimeInfoByTitle(ctx context.Context, title string) (*types.AnimeInfo, error)
}

type AnimeServiceProvider interface {
	Search(query string) ([]types.Result, error)
	FetchEpisodes(id string) ([]types.Episode, error)
}
