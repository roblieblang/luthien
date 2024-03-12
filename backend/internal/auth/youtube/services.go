package youtube

import (
	"context"
	"fmt"
	"net/url"
	"time"
    "log"

	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/utils"
	"google.golang.org/api/youtube/v3"
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

// Completes the initial steps of the authorization code flow with PKCE
func (s *YouTubeService) StartLoginFlow() (string, string, error) {
    log.Printf("Inside StartLoginFlow service")
    sessionID := utils.GenerateSessionID()
    codeVerifier, err := utils.GenerateCodeVerifier(64)
    if err != nil {
        return "","", err
    }

    err = s.YouTubeClient.AppContext.RedisClient.Set(context.Background(), "googleCodeVerifier:" + sessionID, codeVerifier, time.Minute * 10).Err()
    if err != nil {
        return "","", err
    }

    // Request user authorization
    params := url.Values{}
    params.Add("scope", "https://www.googleapis.com/auth/youtube")
    params.Add("access_type", "offline")
    params.Add("prompt", "consent")
    params.Add("include_granted_scopes", "true")
    params.Add("response_type", "code")
    params.Add("state", "state_parameter_passthrough_value") // TODO: use a real value
    params.Add("redirect_uri", s.YouTubeClient.AppContext.EnvConfig.GoogleRedirectURI)
    params.Add("client_id", s.YouTubeClient.AppContext.EnvConfig.GoogleClientID)

    // URL to which the user will be redirected so that they can grant permissions to our application
    authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()

    return authURL, sessionID, nil
} 


// Handles the callback after user has successfully authorized our app on Google's auth page
func (s *YouTubeService) HandleCallback(code, userID, sessionID string) error {
    log.Printf("Inside HandleCallback service")
    payload := url.Values{}
    payload.Set("client_id", s.YouTubeClient.AppContext.EnvConfig.GoogleClientID)
    payload.Set("client_secret", s.YouTubeClient.AppContext.EnvConfig.GoogleClientSecret)
    payload.Set("code", code)
    payload.Set("grant_type", "authorization_code")
    payload.Set("redirect_uri", s.YouTubeClient.AppContext.EnvConfig.GoogleRedirectURI)

    tokenResponse, err := s.YouTubeClient.RequestToken(payload)
    if err != nil {
        return fmt.Errorf("error requesting access token from Google: %v", err)
    }

    if tokenResponse.AccessToken == "" {
        return fmt.Errorf("empty access token")
    }

    params := utils.SetTokenParams{
        TokenKind: "access",
        Party: "google",
        UserID: userID, 
        Token: tokenResponse.AccessToken,
        ExpiresIn: tokenResponse.ExpiresIn,
        AppCtx: *s.YouTubeClient.AppContext,
    }
    err = utils.SetToken(params)
    if err != nil {
        return fmt.Errorf("error storing the access token: %v", err)
    }

    // This is an arbitrary expiry. utils.SetToken() handles refresh token expiration time
    params.TokenKind = "refresh"
    params.Token = tokenResponse.RefreshToken
    log.Printf("Received Google refresh token: '%s'", tokenResponse.RefreshToken)
    params.ExpiresIn = 0
    err = utils.SetToken(params) 
    if err != nil {
        return fmt.Errorf("error storing the Google refresh token: %v", err)
    }
    log.Printf("Google refresh token stored successfully: %s", tokenResponse.RefreshToken)

    // Change user's Google authentication status to `true` 
    updatedAuthStatus := map[string]interface{}{
        "app_metadata": map[string]bool{
            "authenticated_with_google": true,
        },
    }
    if err := s.Auth0Service.UpdateUserMetadata(userID, updatedAuthStatus); err != nil {
        return fmt.Errorf("error updating user metadata: %v", err)
    }
    return nil
}

// Wrapper service function for GetCurrentUserPlaylists client function
func (s *YouTubeService) GetCurrentUserPlaylists(userID string)  (YouTubePlaylistsResponse, error) {
    log.Printf("Inside GetCurrentUserPlaylists service")
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "google", 
        Service: s.YouTubeClient,
        AppCtx: *s.YouTubeClient.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return YouTubePlaylistsResponse{}, err
    }
    return s.YouTubeClient.GetCurrentUserPlaylists(accessToken)
}

// Wrapper service function for GetPlaylistItems client function
func (s *YouTubeService) GetPlaylistItems(userID, playlistID string)  (YouTubePlaylistItemsResponse, error) {
    log.Printf("Inside GetPlaylistItems service")
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "google", 
        Service: s.YouTubeClient,
        AppCtx: *s.YouTubeClient.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    log.Printf("\nGoogle Access Token: '%s'\n", accessToken)
    if err != nil {
        return YouTubePlaylistItemsResponse{}, err
    }
    return s.YouTubeClient.GetPlaylistItems(playlistID, accessToken)
}

// Wrapper service function for CreatePlaylist client function
func (s *YouTubeService) CreatePlaylist(userID string, payload CreatePlaylistPayload) (*youtube.Playlist, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "google", 
        Service: s.YouTubeClient,
        AppCtx: *s.YouTubeClient.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    log.Printf("YouTubeService.CreatePlaylist() access token: '%s'\n", accessToken)
    if err != nil {
        return nil, err
    }
    return s.YouTubeClient.CreatePlaylist(accessToken, payload)
}

// Wrapper service function for AddItemsToPlaylist client function
func (s *YouTubeService) AddItemsToPlaylist(userID string, payload AddItemsToPlaylistPayload) error {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "google", 
        Service: s.YouTubeClient,
        AppCtx: *s.YouTubeClient.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return err
    }
    return s.YouTubeClient.AddItemsToPlaylist(accessToken, payload)
}

// Wrapper service function for SearchVideos client function
func (s *YouTubeService) SearchVideos (userID, artistName, songTitle string) (YouTubeVideoSearchResponse, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "google", 
        Service: s.YouTubeClient,
        AppCtx: *s.YouTubeClient.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return YouTubeVideoSearchResponse{}, err
    }
    query := fmt.Sprintf("%s %s", artistName, songTitle)
    return s.YouTubeClient.SearchVideos(accessToken, query, 1) // maxResults currently hardcoded
}