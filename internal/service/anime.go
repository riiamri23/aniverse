package service

import (
	"time"

	"github.com/Ithilias/anilistgo"
)

type AnimeService struct {
	api *anilistgo.AuthenticatedAPI
}

func NewAnimeService(accessToken string) *AnimeService {
	return &AnimeService{
		api: anilistgo.NewAuthenticatedAPI(accessToken),
	}
}

func (s *AnimeService) GetAnimeByID(id int) (anilistgo.AnilistItem, error) {
	return anilistgo.GetAnilistItemByID(id)
}

func (s *AnimeService) FindAnime(title string, firstEpisodeDate *time.Time, offset int) (anilistgo.AnilistItem, error) {
	return anilistgo.FindAnilistItem(title, firstEpisodeDate, offset)
}

func (s *AnimeService) GetFollowingNames(username string) ([]string, error) {
	return anilistgo.GetFollowingNames(username)
}

func (s *AnimeService) GetUserUpdates(username string, mediaType string, chunk *int, perChunk *int) ([]anilistgo.Update, error) {
	return anilistgo.GetUpdates(username, mediaType, chunk, perChunk)
}

func (s *AnimeService) GetUserProgress(username string, mediaID int) (int, error) {
	return anilistgo.GetProgress(username, mediaID)
}

func (s *AnimeService) UpdateUserProgress(mediaID int, progress int, status string) error {
	return s.api.UpdateProgress(mediaID, progress, status)
}
