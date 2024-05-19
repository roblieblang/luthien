package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien/backend/internal/auth/spotify"
	"github.com/roblieblang/luthien/backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func NewTestSpotifyHandler() *spotify.SpotifyHandler {
	mockSpotifyService := new(MockSpotifyService)
	return spotify.NewSpotifyHandler(mockSpotifyService)
}

func setupRouter(handler *spotify.SpotifyHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/auth/spotify/login", handler.LoginHandler)
	router.POST("/auth/spotify/callback", handler.CallbackHandler)
	router.GET("/spotify/current-profile", handler.GetCurrentUserProfileHandler)
	router.POST("/auth/spotify/logout", handler.LogoutHandler)
	router.GET("/spotify/search-for-track", handler.SearchTracksUsingArtistAndTrackhandler)
	return router
}


func TestSpotifyLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTestSpotifyHandler()
	router := setupRouter(handler)

	mockSpotifyService := handler.SpotifyService.(*MockSpotifyService)
	mockSpotifyService.On("StartLoginFlow").Return("https://example.com/auth", "mockSessionID", nil)

	req, _ := http.NewRequest("GET", "/auth/spotify/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"authURL": "https://example.com/auth", "sessionID": "mockSessionID"}`, w.Body.String())
	mockSpotifyService.AssertExpectations(t)
}

func TestSpotifyCallbackHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTestSpotifyHandler()

	t.Run("valid callback", func(t *testing.T) {
		router := setupRouter(handler)

		mockSpotifyService := handler.SpotifyService.(*MockSpotifyService)
		mockSpotifyService.On("HandleCallback", "authCode", "user123", "session123").Return(nil)

		validBody := `{"code": "authCode", "userID": "user123", "sessionID": "session123"}`
		req, _ := http.NewRequest("POST", "/auth/spotify/callback", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"redirectURL": "http://localhost:5173/"}`, w.Body.String())
		mockSpotifyService.AssertExpectations(t)
	})

	t.Run("invalid callback", func(t *testing.T) {
		router := setupRouter(handler)

		invalidBody := `{"code": "", "userID": "user123", "sessionID": "session123"}`
		req, _ := http.NewRequest("POST", "/auth/spotify/callback", strings.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Missing required fields"}`, w.Body.String())
	})
}

func TestSpotifyGetCurrentUserProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTestSpotifyHandler()

	t.Run("valid request", func(t *testing.T) {
		router := setupRouter(handler)

		mockSpotifyService := handler.SpotifyService.(*MockSpotifyService)
		expectedProfile := spotify.SpotifyUserProfile{
			Country:        "",
			DisplayName:    "Mock User",
			Email:          "",
			ExplicitContent: spotify.ExplicitContent{FilterEnabled: false, FilterLocked: false},
			ExternalUrls:   spotify.ExternalUrls{Spotify: ""},
			Followers:      spotify.Followers{Href: "", Total: 0},
			Href:           "",
			ID:             "user123",
			Images:         nil,
			Product:        "",
			Type:           "",
			URI:            "",
		}
		mockSpotifyService.On("GetCurrentUserProfile", "user123").Return(expectedProfile, nil)

		req, _ := http.NewRequest("GET", "/spotify/current-profile?userID=user123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := `{
			"country": "",
			"display_name": "Mock User",
			"email": "",
			"explicit_content": {
				"filter_enabled": false,
				"filter_locked": false
			},
			"external_urls": {
				"spotify": ""
			},
			"followers": {
				"href": "",
				"total": 0
			},
			"href": "",
			"id": "user123",
			"images": null,
			"product": "",
			"type": "",
			"uri": ""
		}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
		mockSpotifyService.AssertExpectations(t)
	})

	t.Run("missing userID", func(t *testing.T) {
		router := setupRouter(handler)

		req, _ := http.NewRequest("GET", "/spotify/current-profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "userID query parameter is required"}`, w.Body.String())
	})
}



func TestSpotifySearchTracksUsingArtistAndTrackHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTestSpotifyHandler()

	t.Run("valid search", func(t *testing.T) {
		router := setupRouter(handler)
		mockSpotifyService := handler.SpotifyService.(*MockSpotifyService)

		expectedTracks := []utils.UnifiedTrackSearchResult{{ID: "track123", Title: "TrackName", Album: "", Artist: "", Thumbnail: ""}}
		mockSpotifyService.On("SearchTracksUsingArtistAndTrack", "user123", "artist", "track", 20, 0).Return(expectedTracks, nil)

		req, _ := http.NewRequest("GET", "/spotify/search-for-track?userID=user123&artistName=artist&trackTitle=track", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := `[{"id": "track123", "title": "TrackName", "album": "", "artist": "", "thumbnail": ""}]`
		assert.JSONEq(t, expectedResponse, w.Body.String())
		mockSpotifyService.AssertExpectations(t)
	})

	t.Run("missing userID", func(t *testing.T) {
		router := setupRouter(handler)

		req, _ := http.NewRequest("GET", "/spotify/search-for-track?artistName=artist&trackTitle=track", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "userID query parameter is required"}`, w.Body.String())
	})
}

func TestIntegrationSpotifyLoginFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTestSpotifyHandler()
	router := setupRouter(handler)

	mockSpotifyService := handler.SpotifyService.(*MockSpotifyService)
	mockSpotifyService.On("StartLoginFlow").Return("https://example.com/auth", "mockSessionID", nil)
	mockSpotifyService.On("HandleCallback", "authCode", "user123", "session123").Return(nil)

	// Simulate Spotify login flow
	req, _ := http.NewRequest("GET", "/auth/spotify/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "authURL")
	assert.Contains(t, w.Body.String(), "sessionID")

	// Simulate Spotify callback
	validBody := `{"code": "authCode", "userID": "user123", "sessionID": "session123"}`
	req, _ = http.NewRequest("POST", "/auth/spotify/callback", strings.NewReader(validBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"redirectURL": "http://localhost:5173/"}`, w.Body.String())
}
