package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"

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

    err = s.setToken("Access", userID, tokenResponse.AccessToken, tokenResponse.ExpiresIn)
    if err != nil {
        return fmt.Errorf("error storing the access token: %v", err)
    }
    // This is an arbitrary expiry. setToken() handles refresh token expiration time
    err = s.setToken("Refresh", userID, tokenResponse.RefreshToken, 0) 
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

// Called when user clicks "Log Out of Spotify" button on the user interface
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

// Get a token from Redis
func (s *SpotifyService) retrieveToken(tokenKind, userID string) (string, error) {
    token, err := s.AppContext.RedisClient.Get(context.Background(), fmt.Sprintf("spotify%sToken:%s", tokenKind, userID)).Result()
    // Token not found
    if err == redis.Nil {
        return "", nil
    } else if err != nil {
        return "", err
    }
    return token, nil
}

// Store a token in Redis
func (s *SpotifyService) setToken(tokenKind, userID, token string, expiresIn int) error {
    var expiration time.Duration
    if tokenKind == "Access" {
        expiration = time.Duration(expiresIn) * time.Second
    } else if tokenKind == "Refresh" {
        expiration = time.Hour * 720 // one month
    }

    err := s.AppContext.RedisClient.Set(context.Background(), fmt.Sprintf("spotify%sToken:%s", tokenKind, userID), token, expiration).Err()
    if err != nil {
        return fmt.Errorf("error storing the access token: %v", err)
    }
    return nil
}

// Attempts to get a valid access token or sends notice that the user must reauthenticate
func (s *SpotifyService) getValidAccessToken(userID string) (string, error) {
    // Try to get an access token directly from Redis first 
    accessToken, err := s.retrieveToken("Access", userID)
    if err == redis.Nil {
        log.Println("Access token not found in Redis, requesting a new one.")
    } else if err != nil {
        log.Printf("Failed to retrieve Spotify API Access Token: %v", err)
        return "", err
    } else {
        isExpired, err := utils.IsAccessTokenExpired(s.AppContext, "spotifyAccessToken:"+userID)
        if err != nil {
            log.Printf("Failed to check token freshness: %v", err)
            return "", err
        }
        if !isExpired {
            // The token is valid and not expired
            return accessToken, nil
        }
    }
    // If the code reaches here, it means the access token was either not found, expired, or some error occurred
    refreshToken, err := s.retrieveToken("Refresh", userID)
    if err == redis.Nil {
        if err := s.HandleLogout(userID); err != nil {
            log.Printf("Error handling forced logout for user %s: %v", userID, err)
        }
        return "", fmt.Errorf("user must reauthenticate with Spotify")
    } else if err != nil {
        log.Printf("failed to retrieve Spotify API Refresh Token: %v", err)
        return "", err
    } else {
        isExpired, err := utils.IsAccessTokenExpired(s.AppContext, "spotifyRefreshToken:"+userID)
        if err != nil {
            log.Printf("Failed to check token freshness: %v", err)
            return "", err
        }
        // Use valid refresh token to request a new access token from Spotify
        if !isExpired {
            payload := url.Values{}
            payload.Set("grant_type", "refresh_token")
            payload.Set("refresh_token", refreshToken)
            payload.Set("client_id", s.AppContext.EnvConfig.SpotifyClientID)

            tokenResponse, err := s.SpotifyClient.RequestToken(payload)
            if err != nil {
                return "", fmt.Errorf("error requesting access token from Spotify: %v", err)
            }
            if tokenResponse.AccessToken == "" {
                return "", errors.New("empty access token")
            }

            if err := s.setToken("Access", userID, tokenResponse.AccessToken, tokenResponse.ExpiresIn); err != nil {
                return "", fmt.Errorf("error storing access token in Redis: %v", err)
            }
            if err := s.setToken("Refresh", userID, tokenResponse.RefreshToken, 0); err != nil {
                return "", fmt.Errorf("error storing refresh token in Redis: %v", err)
            }
            // Successfully requested a new access token from Spotify using the refresh token
            return tokenResponse.AccessToken, nil
        } else {
            if err := s.HandleLogout(userID); err != nil {
                log.Printf("Error handling forced logout for user %s: %v", userID, err)
            }
            return "", fmt.Errorf("user must reauthenticate with Spotify")
        }
    }
}

// Delete Spotify access and refresh tokens from Redis
func (s *SpotifyService) ClearTokens(userID string) error {
    _, err := s.AppContext.RedisClient.Del(context.Background(), "spotifyAccessToken:" + userID).Result()
    if err != nil {
        return fmt.Errorf("error deleting the access token: %v", err)
    }
    _, err = s.AppContext.RedisClient.Del(context.Background(), "spotifyRefreshToken:" + userID).Result()
    if err != nil {
        return fmt.Errorf("error deleting the refresh token: %v", err)
    }
    return nil
}

// Wrapper service function for GetCurrentUserProfile client function
func (s *SpotifyService) GetCurrentUserProfile(userID string) (SpotifyUserProfile, error) {
    accessToken, err := s.getValidAccessToken(userID)
    if err != nil {
        return SpotifyUserProfile{}, err
    }
    return s.SpotifyClient.GetCurrentUserProfile(accessToken)
}

// Wrapper service function for GetCurrentUserPlaylists client function
func (s *SpotifyService) GetCurrentUserPlaylists(userID string, limit, offset int) (SpotifyPlaylistsResponse, error) {
    accessToken, err := s.getValidAccessToken(userID)
    if err != nil {
        return SpotifyPlaylistsResponse{}, err
    }
    return s.SpotifyClient.GetCurrentUserPlaylists(accessToken, limit, offset)
}

// Wrapper service function for GetPlaylistTracks client function
func (s *SpotifyService) GetPlaylistTracks(userID, playlistID string, limit, offset int) (SpotifyPlaylistTracksResponse, error) {
    accessToken, err := s.getValidAccessToken(userID)
    if err != nil {
        return SpotifyPlaylistTracksResponse{}, err
    }
    return s.SpotifyClient.GetPlaylistTracks(accessToken, playlistID, limit, offset)
}

// Wrapper service function for CreatePlaylist client function
func (s *SpotifyService) CreatePlaylist(userID, spotifyUserID string, payload CreatePlaylistPayload) ([]byte, error) {
    accessToken, err := s.getValidAccessToken(userID)
    if err != nil {
        return nil, err
    }
    return s.SpotifyClient.CreatePlaylist(accessToken, spotifyUserID, payload)
}