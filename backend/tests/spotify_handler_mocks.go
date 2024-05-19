package tests

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/auth/spotify"
	"github.com/roblieblang/luthien/backend/internal/utils"
	"github.com/stretchr/testify/mock"
)

type MockSpotifyService struct {
	mock.Mock
}

func (m *MockSpotifyService) StartLoginFlow() (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockSpotifyService) HandleCallback(code, userID, sessionID string) error {
	args := m.Called(code, userID, sessionID)
	return args.Error(0)
}

func (m *MockSpotifyService) GetCurrentUserProfile(userID string) (spotify.SpotifyUserProfile, error) {
	args := m.Called(userID)
	return args.Get(0).(spotify.SpotifyUserProfile), args.Error(1)
}

func (m *MockSpotifyService) GetCurrentUserPlaylists(userID string, offset int) (spotify.SpotifyPlaylistsResponse, error) {
	args := m.Called(userID, offset)
	return args.Get(0).(spotify.SpotifyPlaylistsResponse), args.Error(1)
}

func (m *MockSpotifyService) GetPlaylistTracks(userID, playlistID string) (spotify.SpotifyPlaylistTracksResponse, error) {
	args := m.Called(userID, playlistID)
	return args.Get(0).(spotify.SpotifyPlaylistTracksResponse), args.Error(1)
}

func (m *MockSpotifyService) CreatePlaylist(userID, spotifyUserID string, payload spotify.CreatePlaylistPayload) (string, error) {
	args := m.Called(userID, spotifyUserID, payload)
	return args.String(0), args.Error(1)
}

func (m *MockSpotifyService) AddItemsToPlaylist(userID, playlistID string, payload spotify.AddItemsToPlaylistPayload) error {
	args := m.Called(userID, playlistID, payload)
	return args.Error(0)
}

func (m *MockSpotifyService) SearchTracksUsingArtistAndTrack(userID, artistName, trackTitle string, limit, offset int) ([]utils.UnifiedTrackSearchResult, error) {
	args := m.Called(userID, artistName, trackTitle, limit, offset)
	return args.Get(0).([]utils.UnifiedTrackSearchResult), args.Error(1)
}

func (m *MockSpotifyService) SearchTracksUsingVideoTitle(userID, videoTitle string) ([]utils.UnifiedTrackSearchResult, error) {
	args := m.Called(userID, videoTitle)
	return args.Get(0).([]utils.UnifiedTrackSearchResult), args.Error(1)
}

func (m *MockSpotifyService) DeletePlaylist(userID, playlistID string) error {
	args := m.Called(userID, playlistID)
	return args.Error(0)
}

func (m *MockSpotifyService) GetAuth0Service() *auth0.Auth0Service {
	args := m.Called()
	return args.Get(0).(*auth0.Auth0Service)
}

func (m *MockSpotifyService) GetAppContext() *utils.AppContext {
	args := m.Called()
	return args.Get(0).(*utils.AppContext)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	var args []interface{}
	for _, key := range keys {
		args = append(args, key)
	}
	return redis.NewIntCmd(ctx, args...)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx, key, value, expiration)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringCmd(ctx, key)
}

func (m *MockRedisClient) TTL(ctx context.Context, key time.Duration) *redis.DurationCmd {
	return redis.NewDurationCmd(ctx, key)
}
