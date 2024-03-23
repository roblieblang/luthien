package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type UnifiedTrackSearchResult struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Artist      string `json:"artist"`
    Album       string `json:"album"`
    Thumbnail   string `json:"thumbnail"`
}

// Returns a high-entropy random string which will be used as a code verifier after being hashed
func GenerateCodeVerifier(length int) (string, error) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.-~"
    result := make([]byte, length)
    for i := range result {
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
        if err != nil {
            return "", err
        }
        result[i] = chars[num.Int64()]
    }
    return string(result), nil
}

// Returns a SHA256 hashed and base64 encoded string which can be used as a code verifier
func SHA256Hash(data string) string {
    hash := sha256.Sum256([]byte(data))
    return base64.RawURLEncoding.EncodeToString(hash[:])
}

// Returns a randomized session ID
func GenerateSessionID() string {
    b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

// Capitalize first letter of a string
func capitalizeFirstLetter(s string) string {
    if s == "" {
        return ""
    }
    return strings.ToUpper(string(s[0])) + s[1:]
}

type SetTokenParams struct {
    TokenKind   string
    Party       string
    UserID      string
    Token       string
    ExpiresIn   int
    AppCtx      AppContext
}

// Stores a token in Redis
func SetToken(params SetTokenParams) error {
    log.Printf("\nSetToken() params: %v\n", params)
    party := strings.ToLower(params.Party)
    tokenKind := capitalizeFirstLetter(params.TokenKind)

    if params.Token == "" {
        return fmt.Errorf("empty %s token passed into SetToken()", params.TokenKind)
    }

    var expiration time.Duration
    if tokenKind == "Access" {
        expiration = time.Duration(params.ExpiresIn) * time.Second
    } else if tokenKind  == "Refresh" {
        expiration = time.Hour * 720 // one month
    }
    log.Printf("Setting %s %s token for user %s with value: %s", party, tokenKind, params.UserID, params.Token)
    err := params.AppCtx.RedisClient.Set(context.Background(), fmt.Sprintf("%s%sToken:%s", party, tokenKind, params.UserID), params.Token, expiration).Err()
    if err != nil {
        return fmt.Errorf("error storing the access token: %v", err)
    }
    return nil
}

type ClearTokensParams struct {
    Party       string
    UserID      string
    AppCtx  AppContext
}

// Delete access and refresh tokens from Redis
func ClearTokens(params ClearTokensParams) error {
    party := strings.ToLower(params.Party)

    _, err := params.AppCtx.RedisClient.Del(context.Background(), fmt.Sprintf("%sAccessToken:%s", party, params.UserID)).Result()
    if err != nil {
        return fmt.Errorf("error deleting the access token: %v", err)
    }
    _, err = params.AppCtx.RedisClient.Del(context.Background(), fmt.Sprintf("%sRefreshToken:%s", party, params.UserID)).Result()
    if err != nil {
        return fmt.Errorf("error deleting the refresh token: %v", err)
    }
    return nil
}

type RetrieveTokenParams struct {
    Party       string
    TokenKind   string
    UserID      string
    AppCtx  AppContext
}

// Get a token from Redis
func RetrieveToken(params RetrieveTokenParams) (string, error) {
    party := strings.ToLower(params.Party)
    tokenKind := capitalizeFirstLetter(params.TokenKind)

    token, err := params.AppCtx.RedisClient.Get(context.Background(), fmt.Sprintf("%s%sToken:%s", party, tokenKind, params.UserID)).Result()
    // Token not found
    if err == redis.Nil {
        log.Printf("%s %s token not found for user %s", party, tokenKind, params.UserID)
        return "", err
    } else if err != nil {
        log.Printf("error whiole retrieving token: %v", err)
        return "", err
    } else if token == "" {
        // Token is found but its value is empty
        log.Printf("%s %s token found with empty value for user %s", party, tokenKind, params.UserID)
        return "", nil
    }
    log.Printf("Retrieved %s %s token for user %s with value: %s", party, tokenKind, params.UserID, token)
    return token, nil
}

// Checks if the access token stored in Redis identified by `tokenName` argument has expired
func IsAccessTokenExpired(appCtx AppContext, tokenName, token string) (bool, error) {
    if token == "" {
        log.Printf("Token (%s: %s) is empty\n", tokenName, token)
        return true, nil
    }

    accessTokenTTL, err := appCtx.RedisClient.TTL(context.Background(), tokenName).Result()
    if err != nil {
        log.Printf("error retrieving token (%s: %s) time to live: %v\n", tokenName, token, err)
        return true, err
    }

    // Token has a TTL and is soon to expire
    if accessTokenTTL > 0 && accessTokenTTL < 5 * time.Minute {
        log.Printf("Token (%s: %s) is soon to expire\n", tokenName, token)
        return true, nil
    }

    if accessTokenTTL == -2 {
        log.Printf("Token (%s: %s) does not exist\n", tokenName, token)
        return true, nil
    }

    if accessTokenTTL == -1 {
        log.Printf("Token (%s: %s) does not have expiry set\n", tokenName, token)
        return false, nil
    }
    
    // Token is not expired
    return false, nil
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    Scope        string `json:"scope"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"` // Make RefreshToken optional
}

type TokenService interface {
    RequestToken(payload url.Values) (TokenResponse, error)
}   

// Must specify spotify or google as "Party"
type GetValidAccessTokenParams struct {
    UserID      string
    Party       string
    Service     TokenService
    AppCtx      AppContext
    Updater     UserMetadataUpdater 
}

func getClientID(party string, appCtx AppContext) (string, error) {
    party = strings.ToLower(party)
    if party == "google" {
        return appCtx.EnvConfig.GoogleClientID, nil
    } else if party == "spotify" {
        return appCtx.EnvConfig.SpotifyClientID, nil
    }
    return "", errors.New("getClientID() only works for Spotify and Google services")
}

// Attempts to get a valid access token or sends notice that the user must reauthenticate
func GetValidAccessToken(params GetValidAccessTokenParams) (string, error) {
    // Try to get an access token directly from Redis first
    accessToken, err := RetrieveToken(RetrieveTokenParams{
        Party: params.Party, 
        TokenKind: "access", 
        UserID: params.UserID, 
        AppCtx: params.AppCtx,
    })
    
    if err != nil && err != redis.Nil {
        log.Printf("Unexpected error occurred while attempting to retrieve an access token: %v", err)
        return "", err
    } 

    isExpired, err := IsAccessTokenExpired(params.AppCtx, fmt.Sprintf("%sAccessToken:%s", params.Party, params.UserID), accessToken)
    if err != nil {
        log.Printf("Error checking if access token is expired: %v", err)
        return "", err
    }
    if !isExpired {
        // The token is valid and not expired
        log.Printf("Retrieved valid access token for user %s with value: %s", params.UserID, accessToken)
        return accessToken, nil
    }

    // If the access token is expired or not found, attempt to use the refresh token
    refreshToken, err := RetrieveToken(RetrieveTokenParams{
        Party: params.Party,
        TokenKind: "refresh",
        UserID: params.UserID,
        AppCtx: params.AppCtx,
    })

    if refreshToken == "" || err == redis.Nil {
        // Refresh token is missing or empty; force reauthentication
        log.Printf("Absent refresh token. Logging out user %s", params.UserID)
        if err := HandleLogout(params.Updater, ClearTokensParams{
            Party: params.Party,
            UserID: params.UserID,
            AppCtx: params.AppCtx,
        }); err != nil {
            log.Printf("Error handling forced logout for user %s: %v", params.UserID, err)
            return "", fmt.Errorf("error forcing logout for user %s: %v", params.UserID, err)
        }
        return "", fmt.Errorf("reauthentication required with %s", params.Party)
    } else if err != nil {
        log.Printf("Error retrieving refresh token: %v", err)
        return "", err
    }

    isRefreshExpired, err := IsAccessTokenExpired(params.AppCtx, fmt.Sprintf("%sRefreshToken:%s", params.Party, params.UserID), refreshToken)
    log.Printf("Refresh Token: '%s'\nisExpired?: %v", fmt.Sprintf("%sRefreshToken:%s", params.Party, params.UserID), isRefreshExpired)
    if err != nil {
        log.Printf("Error checking if refresh token is expired: %v", err)
        return "", err
    }

    // Use valid refresh token to request a new access token from <party>
    if !isRefreshExpired {
        log.Printf("Refreshing %s access token for user %s with refresh token: '%s'", params.Party, params.UserID, refreshToken)

        payload := url.Values{}
        payload.Set("grant_type", "refresh_token")
        payload.Set("refresh_token", refreshToken)
        clientID, err := getClientID(params.Party, params.AppCtx)
        if err != nil {
            log.Printf("Error getting client ID for %s: %v", params.Party, err)
            return "", err
        }
        payload.Set("client_id", clientID)
        if params.Party == "google" || params.Party == "Google" {
            payload.Set("client_secret", params.AppCtx.EnvConfig.GoogleClientSecret)
        }

        log.Printf("%s Access token request payload: %v", params.Party, payload)

        tokenResponse, err := params.Service.RequestToken(payload)
        if err != nil {
            log.Printf("Error refreshing token for user %s with %s, forcing reauthentication: %v", params.UserID, params.Party, err)
            if logoutErr := HandleLogout(params.Updater, ClearTokensParams{
                Party: params.Party,
                UserID: params.UserID,
                AppCtx: params.AppCtx,
            }); logoutErr != nil {
                log.Printf("Error handling forced logout/reauthentication for user %s: %v", params.UserID, logoutErr)
                return "", fmt.Errorf("error forcing logout/reauthentication for user %s: %v", params.UserID, logoutErr)
            }
            return "", fmt.Errorf("reauthentication required with %s for user %s", params.Party, params.UserID)
        }

        log.Printf("Received token response from %s for user %s: AccessToken=%s, ExpiresIn=%d, RefreshToken=%s", params.Party, params.UserID, tokenResponse.AccessToken, tokenResponse.ExpiresIn, tokenResponse.RefreshToken)

        if tokenResponse.AccessToken == "" {
            log.Printf("Refreshed access token came in empty for user %s with %s, forcing reauthentication: %v", params.UserID, params.Party, err)
            if logoutErr := HandleLogout(params.Updater, ClearTokensParams{
                Party: params.Party,
                UserID: params.UserID,
                AppCtx: params.AppCtx,
            }); logoutErr != nil {
                log.Printf("Error handling forced logout/reauthentication for user %s: %v", params.UserID, logoutErr)
                return "", fmt.Errorf("error forcing logout/reauthentication for user %s: %v", params.UserID, logoutErr)
            }
            return "", fmt.Errorf("reauthentication required with %s for user %s", params.Party, params.UserID)
        }

        // Store the access token
        if err := SetToken(SetTokenParams{
            TokenKind: "access",
            Party: params.Party,
            UserID: params.UserID, 
            Token: tokenResponse.AccessToken,
            ExpiresIn: tokenResponse.ExpiresIn,
            AppCtx: params.AppCtx,
        }); err != nil {
            return "", fmt.Errorf("error storing new access token: %v", err)
        }

        // Store the new refresh token if present
        if tokenResponse.RefreshToken != "" {
            if err := SetToken(SetTokenParams{
                TokenKind: "refresh",
                Party: params.Party,
                UserID: params.UserID,
                Token: tokenResponse.RefreshToken,
                ExpiresIn: 0,
                AppCtx: params.AppCtx,
            }); err != nil {
                return "", fmt.Errorf("error storing new refresh token: %v", err)
            }
        }
        // Successfully requested a new access token from <party> using the refresh token
        return tokenResponse.AccessToken, nil
    } 

    // Refresh token is expired, so user must reauthenticate
    if err := HandleLogout(params.Updater, ClearTokensParams{
        Party: params.Party,
        UserID: params.UserID,
        AppCtx: params.AppCtx,
    }); err != nil {
        log.Printf("Error handling forced logout for user %s: %v", params.UserID, err)
        return "", fmt.Errorf("error forcing logout for user %s: %v", params.UserID, err)
    }
    return "", fmt.Errorf("reauthentication required with %s", params.Party)
}

// To mitigate circular dependency between auth0 and utils packages
type UserMetadataUpdater interface {
    UpdateUserMetadata(userID string, updatedAuthStatus map[string]interface{}) error
}

// Called when user clicks "Log Out of <party>" button on the user interface
func HandleLogout(updater UserMetadataUpdater, params ClearTokensParams) error {
    log.Printf("Inside HandleLogout util")
    if err := ClearTokens(params); err != nil {
        return fmt.Errorf("error clearing %s tokens from Redis: %v", params.Party, err)
    }

    party := strings.ToLower(params.Party)

    updatedAuthStatus := map[string]interface{}{
        "app_metadata": map[string]bool{
            fmt.Sprintf("authenticated_with_%s", party): false,
        },
    }
    if err := updater.UpdateUserMetadata(params.UserID, updatedAuthStatus); err != nil {
        return fmt.Errorf("error updating user metadata: %v", err)
    }
    return nil
}