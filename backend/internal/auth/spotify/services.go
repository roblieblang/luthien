package spotify

import (
	"context"
	"errors"
	"fmt"
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

    scope := "user-read-private user-read-email"
    params := url.Values{}
    params.Add("client_id", s.AppContext.EnvConfig.SpotifyClientID)
    params.Add("response_type", "code")
    params.Add("redirect_uri", s.AppContext.EnvConfig.SpotifyRedirectURI)
    params.Add("scope", scope)
    params.Add("code_challenge_method", "S256")
    params.Add("code_challenge", codeChallenge)

    authURL := "https://accounts.spotify.com/authorize?" + params.Encode()

    return authURL, sessionID, nil
} 

func (s *SpotifyService) HandleCallback(code, userID, sessionID string) error {
    codeVerifier, err := s.AppContext.RedisClient.Get(context.Background(), "spotifyCodeVerifier:"+sessionID).Result()
    if err != nil {
        return fmt.Errorf("error retrieving the code verifier: %v", err)
    }

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

    expiresIn := time.Duration(tokenResponse.ExpiresIn) * time.Second
    if tokenResponse.AccessToken == "" {
        return errors.New("empty access token")
    }

    err = s.AppContext.RedisClient.Set(context.Background(), "spotifyAccessToken:"+userID, tokenResponse.AccessToken, expiresIn).Err()
    if err != nil {
        return fmt.Errorf("error storing the access token: %v", err)
    }
    err = s.AppContext.RedisClient.Set(context.Background(), "spotifyRefreshToken:"+userID, tokenResponse.RefreshToken, expiresIn).Err()
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

func (s *SpotifyService) HandleLogout(userID string) error {
    if err := s.ClearTokens(userID); err != nil {
        return fmt.Errorf("error clearing tokens from Redis: %v", err)
    }

    // Change user's Spotify authentication status to `false` 
    updatedAuthStatus := map[string]interface{}{
        "app_metadata": map[string]bool{
            "authenticated_with_spotify": false,
        },
    }
    if err := s.Auth0Service.UpdateUserMetadata(userID, updatedAuthStatus); err != nil {
        return fmt.Errorf("error updating user metadata: %v", err)
    }
    return nil
}

// Delete the Spotify access and refresh tokens from Redis
func (s *SpotifyService) ClearTokens(userID string) error {
    _, err := s.AppContext.RedisClient.Del(context.Background(), "spotifyAccessToken:" + userID).Result()
    if err != nil {
        return fmt.Errorf("There was an issue deleting the access token: %v\n", err)
    }
    _, err = s.AppContext.RedisClient.Del(context.Background(), "spotifyRefreshToken:" + userID).Result()
    if err != nil {
        return fmt.Errorf("There was an issue deleting the refresh token: %v\n", err)
    }
    return nil
}