package youtube

import (
	"fmt"

	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
)

type YouTubeService struct {
	YouTubeClient *YouTubeClient
	Auth0Service  *auth0.Auth0Service
}

func NewYouTubeService(youTubeClient *YouTubeClient, auth0Service *auth0.Auth0Service) *YouTubeService {
	return &YouTubeService{
		YouTubeClient: youTubeClient,
		Auth0Service:  auth0Service,
	}
}

// Wrapper service function for GetCurrentUserPlaylists client function
func (s *YouTubeService) GetCurrentUserPlaylists(userID string)  (YouTubePlaylistsResponse, error) {
	accessToken, err := s.Auth0Service.RetrieveGoogleToken(userID)
	if err != nil {
		return YouTubePlaylistsResponse{}, fmt.Errorf("error retrieving google access token: %v", err)
	}
	return s.YouTubeClient.GetCurrentUserPlaylists(accessToken)
}