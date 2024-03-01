package auth0

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

type Auth0Service struct {
    Auth0Client *Auth0Client
    AppContext *utils.AppContext
}

func NewAuth0Service(auth0Client *Auth0Client, appCtx *utils.AppContext) *Auth0Service {
    return &Auth0Service{
        Auth0Client: auth0Client,
        AppContext: appCtx,
    }
}

// Stores an Auth0 Management API access token in Redis
func (s *Auth0Service) storeAuth0Token(tokenResponse utils.TokenResponse) error{
    err := s.AppContext.RedisClient.Set(context.Background(), "auth0ManagementAPIAccessToken", tokenResponse.AccessToken, time.Duration(tokenResponse.ExpiresIn) * time.Second).Err()
    if err != nil {
        log.Printf("There was an issue storing the Auth0 Management API Access Token: %v", err)
        return err
    }
    return nil
}

// Retrieves an existing Auth0 Management API access token from Redis
func (s *Auth0Service) retrieveAuth0Token() (string, error){
	accessToken, err := s.AppContext.RedisClient.Get(context.Background(), "auth0ManagementAPIAccessToken").Result()
    // Token not found, not an error
    if err == redis.Nil {
        return "", nil
    } else if err != nil {
        log.Printf("Failed to retrieve Auth0 Management API Access Token: %v", err)
        return "", err
    }
    return accessToken, nil
}

// Stores a Google APi access token in Redis
func (s *Auth0Service) storeGoogleToken(userID, googleToken string, expiresIn int) error {
    err := s.AppContext.RedisClient.Set(context.Background(), "googleAPIAccessToken:"+userID, googleToken, time.Duration(expiresIn) * time.Second).Err()
    if err != nil {
        log.Printf("There was an issue storing the Google API Access Token: %v", err)
        return err
    }
    return nil
}

// Retrieves an existing Google API access token from Redis
func (s *Auth0Service) RetrieveGoogleToken(userID string) (string, error){
	accessToken, err := s.AppContext.RedisClient.Get(context.Background(), "googleAPIAccessToken:"+userID).Result()
    // Token not found, not an error
    if err == redis.Nil {
        return "", nil
    } else if err != nil {
        log.Printf("Failed to retrieve Google API Access Token: %v", err)
        return "", err
    }
    return accessToken, nil
}

// Helper function for getting a valid access token
func (s *Auth0Service) getValidAccessToken() (string, error) {
    accessToken, err := s.retrieveAuth0Token()
    // Check for redis.Nil to determine if the key was simply not found i.e. not an actual error
    if err == redis.Nil {
        log.Println("Access token not found in Redis, requesting a new one.")
    } else if err != nil {
        // Handle actual errors from Redis
        log.Printf("Failed to retrieve Auth0 Management API Access Token: %v", err)
        return "", err
    } else {
        // If there was no error and the token exists, check for expiration
        isExpired, err := utils.IsAccessTokenExpired(*s.AppContext, "auth0ManagementAPIAccessToken", accessToken)
        if err != nil {
            log.Printf("Failed to check token freshness: %v", err)
            return "", err
        }
        if !isExpired && accessToken != "" {
            // The token is valid and not expired
            return accessToken, nil
        }
    }

    // If the code reaches here, it means the token was either not found, expired, or another error occurred
    // Attempt to request a new token.
    tokenResponse, err := s.Auth0Client.RequestToken()
    if err != nil {
        log.Printf("Failed to request new Auth0 Management API Access Token: %v", err)
        return "", err
    }
    err = s.storeAuth0Token(tokenResponse)
    if err != nil {
        log.Printf("Failed to store new Auth0 Management API Access Token: %v", err)
        return "", err
    }
    return tokenResponse.AccessToken, nil
}

// Wrapper service function for GetUserMetadata client function that extracts and stores Google access token from response
func (s *Auth0Service) GetUserMetadata(userID string) (Auth0UserMetadata, error) {
    accessToken, err := s.getValidAccessToken()
    if err != nil {
        return Auth0UserMetadata{}, err
    }
    userMetadata, err := s.Auth0Client.GetUserMetadata(accessToken, userID)
    if err != nil {
        return Auth0UserMetadata{}, err
    }

    return userMetadata, nil
}

// Wrapper service function for UpdateUserMetadata client function
func (s *Auth0Service) UpdateUserMetadata(userID string, updates map[string]interface{}) error {
    accessToken, err := s.getValidAccessToken()
    if err != nil {
        return err
    }
    return s.Auth0Client.UpdateUserMetadata(accessToken, userID, updates)
}