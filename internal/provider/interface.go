package provider

import (
	"aniverse/internal/domain/types"
	"context"
)

type InformationProvider interface {
	GetAnimeInfoByID(ctx context.Context, mediaID int) (*types.AnimeInfo, error)
	GetAnimeInfoByTitle(ctx context.Context, title string) (*types.AnimeInfo, error)
}

type AnimeServiceProvider interface {
	Search(query string) ([]types.Result, error)
	Episodes(id string) ([]types.Episode, error)
	Sources(id, subType, server string) (*types.Source, error)
}
