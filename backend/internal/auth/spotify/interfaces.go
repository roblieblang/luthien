package spotify

import (
	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

type SpotifyServiceInterface interface {
	StartLoginFlow() (string, string, error)
	HandleCallback(code, userID, sessionID string) error
	GetCurrentUserProfile(userID string) (SpotifyUserProfile, error)
	GetCurrentUserPlaylists(userID string, offset int) (SpotifyPlaylistsResponse, error)
	GetPlaylistTracks(userID, playlistID string) (SpotifyPlaylistTracksResponse, error)
	CreatePlaylist(userID, spotifyUserID string, payload CreatePlaylistPayload) (string, error)
	AddItemsToPlaylist(userID, playlistID string, payload AddItemsToPlaylistPayload) error
	SearchTracksUsingArtistAndTrack(userID, artistName, trackTitle string, limit, offset int) ([]utils.UnifiedTrackSearchResult, error)
	SearchTracksUsingVideoTitle(userID, videoTitle string) ([]utils.UnifiedTrackSearchResult, error)
	DeletePlaylist(userID, playlistID string) error
	GetAuth0Service() *auth0.Auth0Service
	GetAppContext() *utils.AppContext
}
