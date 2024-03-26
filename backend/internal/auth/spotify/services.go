package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

// Calls our spotify client
// Applies our own application's rules (business logic) to the data received
// Interacts with other parts of our application like Mongo, Redis, other services

type SpotifyService struct {
    SpotifyClient *SpotifyClient
    Auth0Service *auth0.Auth0Service
    AppContext    *utils.AppContext
}

func NewSpotifyService(spotifyClient *SpotifyClient, auth0Service *auth0.Auth0Service, appContext *utils.AppContext) *SpotifyService {
    return &SpotifyService{
        SpotifyClient: spotifyClient,
        Auth0Service: auth0Service,
        AppContext: appContext,
    }
}

// Completes the initial steps of the authorization code flow with PKCE
func (s *SpotifyService) StartLoginFlow() (string, string, error) {
    sessionID := utils.GenerateSessionID()
    codeVerifier, err := utils.GenerateCodeVerifier(64)
    if err != nil {
        return "","", err
    }

    err = s.AppContext.RedisClient.Set(context.Background(), "spotifyCodeVerifier:" + sessionID, codeVerifier, time.Minute * 10).Err()
    if err != nil {
        return "","", err
    }

    codeChallenge := utils.SHA256Hash(codeVerifier)

    // Request user authorization
    scope := "user-read-private user-read-email playlist-read-private playlist-read-collaborative playlist-modify-public playlist-modify-private"
    params := url.Values{}
    params.Add("client_id", s.AppContext.EnvConfig.SpotifyClientID)
    params.Add("response_type", "code")
    params.Add("redirect_uri", s.AppContext.EnvConfig.SpotifyRedirectURI)
    params.Add("scope", scope)
    params.Add("code_challenge_method", "S256")
    params.Add("code_challenge", codeChallenge)

    // URL to which the user will be redirected so that they can grant permissions to our application
    authURL := "https://accounts.spotify.com/authorize?" + params.Encode()

    return authURL, sessionID, nil
} 

// Handles the callback after user has successfully authorized our app on Spotify's auth page
func (s *SpotifyService) HandleCallback(code, userID, sessionID string) error {
    codeVerifier, err := s.AppContext.RedisClient.Get(context.Background(), "spotifyCodeVerifier:"+sessionID).Result()
    if err != nil {
        return fmt.Errorf("error retrieving the code verifier: %v", err)
    }

    // Request an access token
    payload := url.Values{}
    payload.Set("grant_type", "authorization_code")
    payload.Set("code", code)
    payload.Set("redirect_uri", s.AppContext.EnvConfig.SpotifyRedirectURI)
    payload.Set("client_id", s.AppContext.EnvConfig.SpotifyClientID)
    payload.Set("code_verifier", codeVerifier)

    tokenResponse, err := s.SpotifyClient.RequestToken(payload)
    if err != nil {
        return fmt.Errorf("error requesting access token from Spotify: %v", err)
    }

    if tokenResponse.AccessToken == "" {
        return errors.New("empty access token")
    }

    // Store the access token
    params := utils.SetTokenParams{
        TokenKind: "access",
        Party: "spotify",
        UserID: userID, 
        Token: tokenResponse.AccessToken,
        ExpiresIn: tokenResponse.ExpiresIn,
        AppCtx: *s.AppContext,
    }
    err = utils.SetToken(params)
    if err != nil {
        return fmt.Errorf("error storing the access token: %v", err)
    }

    // Now store the access token
    params.TokenKind = "refresh"
    params.Token = tokenResponse.RefreshToken
    // This is an arbitrary expiry. SetToken() handles refresh token expiration time
    params.ExpiresIn = 0
    err = utils.SetToken(params)
    if err != nil {
        return fmt.Errorf("error storing the refresh token: %v", err)
    }

    // Change user's Spotify authentication status to `true` 
    updatedAuthStatus := map[string]interface{}{
        "app_metadata": map[string]bool{
            "authenticated_with_spotify": true,
        },
    }
    if err := s.Auth0Service.UpdateUserMetadata(userID, updatedAuthStatus); err != nil {
        return fmt.Errorf("error updating user metadata: %v", err)
    }
    return nil
}

// Wrapper service function for GetCurrentUserProfile client function
func (s *SpotifyService) GetCurrentUserProfile(userID string) (SpotifyUserProfile, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        log.Printf("error getting a valid spotify access token: %v", err)
        return SpotifyUserProfile{}, err
    }
    return s.SpotifyClient.GetCurrentUserProfile(accessToken)
}

// Wrapper service function for GetCurrentUserPlaylists client function
func (s *SpotifyService) GetCurrentUserPlaylists(userID string, offset int) (SpotifyPlaylistsResponse, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return SpotifyPlaylistsResponse{}, err
    }
    return s.SpotifyClient.GetCurrentUserPlaylists(accessToken, offset)
}

// Wrapper service function for GetPlaylistTracks client function
func (s *SpotifyService) GetPlaylistTracks(userID, playlistID string) (SpotifyPlaylistTracksResponse, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return SpotifyPlaylistTracksResponse{}, err
    }
    return s.SpotifyClient.GetPlaylistTracks(accessToken, playlistID)
}

// Wrapper service function for CreatePlaylist client function
func (s *SpotifyService) CreatePlaylist(userID, spotifyUserID string, payload CreatePlaylistPayload) (string, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return "", err
    }
    return s.SpotifyClient.CreatePlaylist(accessToken, spotifyUserID, payload)
}

// Wrapper service function for AddItemsToPlaylist client function
func (s *SpotifyService) AddItemsToPlaylist(userID, playlistID string, payload AddItemsToPlaylistPayload) error {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        return err
    }
    return s.SpotifyClient.AddItemsToPlaylist(accessToken, playlistID, payload)
}

// Wrapper service function for SearchTracksUsingArtistAndTrack client function
func (s *SpotifyService) SearchTracksUsingArtistAndTrack(userID, artistName, trackTitle string, limit, offset int) ([]utils.UnifiedTrackSearchResult, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        log.Printf("Error getting valid access token: %v", err)
        return nil, err
    }
    return s.SpotifyClient.SearchTracksUsingArtistAndTrack(accessToken, artistName, trackTitle, limit, offset)
}

// Wrapper service function for SearchTracksUsingVideoTitle client function
func (s *SpotifyService) SearchTracksUsingVideoTitle(userID, videoTitle string) ([]utils.UnifiedTrackSearchResult, error) {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        log.Printf("Error getting valid access token: %v", err)
        return nil, err
    }
    return s.SpotifyClient.SearchTracksUsingVideoTitle(accessToken, videoTitle)
}

// Wrapper service function for DeletePlaylist client function
func (s *SpotifyService) DeletePlaylist(userID, playlistID string) error {
    params := utils.GetValidAccessTokenParams{
        UserID: userID, 
        Party: "spotify", 
        Service: s.SpotifyClient,
        AppCtx: *s.AppContext,
        Updater: s.Auth0Service,
    }
    accessToken, err := utils.GetValidAccessToken(params)
    if err != nil {
        log.Printf("Error getting valid access token: %v", err)
        return err
    }

    return s.SpotifyClient.DeletePlaylist(accessToken, playlistID)    
}