package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
    "github.com/redis/go-redis/v9"  
    "fmt"
    "net/url"
	"io"
    "errors"
	"log"
	"math/big"
	"strings"
	"time"
)

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
    AppCtx  AppContext
}

// Stores a token in Redis
func SetToken(params SetTokenParams) error {
    fmt.Printf("\nSetToken() params: %v\n", params)
    party := strings.ToLower(params.Party)
    tokenKind := capitalizeFirstLetter(params.TokenKind)

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
        return "", nil
    } else if err != nil {
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
        return true, nil
    }

    accessTokenTTL, err := appCtx.RedisClient.TTL(context.Background(), tokenName).Result()
    if err != nil {
        log.Printf("There was an issue retrieving the access token time to live: %v\n", err)
        return true, err
    }

    // Token has a TTL and is soon to expire
    if accessTokenTTL > 0 && accessTokenTTL < 5 * time.Minute {
        return true, nil
    }

    if accessTokenTTL == -2 {
        log.Printf("Token key does not exist\n")
        return true, nil
    }

    if accessTokenTTL == -1 {
        log.Printf("Token does not have an expiry set\n")
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
    AppCtx  AppContext
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
    tokenParams := RetrieveTokenParams{
        Party: params.Party, 
        TokenKind: "access", 
        UserID: params.UserID, 
        AppCtx: params.AppCtx,
    } 
    accessToken, err := RetrieveToken(tokenParams)
    clearTokenParams := ClearTokensParams{
        Party: params.Party,
        UserID: params.UserID,
        AppCtx: params.AppCtx,
    }
    if err != nil {
        log.Printf("error occurred while attempting to retrieve an access token: %v", err)
        return "", err
    } else {
        isExpired, err := IsAccessTokenExpired(params.AppCtx, fmt.Sprintf("%sAccessToken:%s", params.Party, params.UserID), accessToken)
        if err != nil {
            log.Printf("Failed to check token freshness: %v", err)
            return "", err
        }
        if !isExpired {
            // The token is valid and not expired
            log.Printf("Retrieved valid access token for user %s with value: %s", params.UserID, accessToken)
            return accessToken, nil
        }
    }
    // If the code reaches here, it means the access token was either not found or expired
    tokenParams.TokenKind = "refresh"
    fmt.Printf("\nREFRESH TOKEN PARAMS: %v\n", tokenParams)
    refreshToken, err := RetrieveToken(tokenParams)
    fmt.Printf("\nREFRESH TOKEN: %s\n", refreshToken)

    if err == redis.Nil {
        // Empty refresh token means that the user's authentication session has expired and they must now reauthenticate
        log.Printf("Absent refresh token. Logging out user %s", params.UserID)
        if err := HandleLogout(params.Updater, clearTokenParams); err != nil {
            log.Printf("Error handling forced logout for user %s: %v", params.UserID, err)
        }
        return "", fmt.Errorf("reauthentication required with %s", params.Party)
    } else if err != nil {
        log.Printf("failed to retrieve %s API Refresh Token: %v", params.Party, err)
        return "", err
    } else {
        isExpired, err := IsAccessTokenExpired(params.AppCtx, fmt.Sprintf("%sRefreshToken:%s", params.Party, params.UserID), refreshToken)
        if err != nil {
            log.Printf("Failed to check token freshness: %v", err)
            return "", err
        }
        // Use valid refresh token to request a new access token from <party>
        if !isExpired {
            log.Printf("Refreshing access token for user %s with refresh token: %s", params.UserID, refreshToken)
            payload := url.Values{}
            payload.Set("grant_type", "refresh_token")
            payload.Set("refresh_token", refreshToken)
            clientID, err := getClientID(params.Party, params.AppCtx)
            if err != nil {
                return "", err
            }
            payload.Set("client_id", clientID)
            if params.Party == "google" || params.Party == "Google" {
                payload.Set("client_secret", params.AppCtx.EnvConfig.GoogleClientSecret)
            }

            fmt.Printf("%s Access token request payload: %v", params.Party, payload)

            tokenResponse, err := params.Service.RequestToken(payload)
            fmt.Printf("\nTOKEN RESPONSE: %v\n", tokenResponse)
            if err != nil {
                return "", fmt.Errorf("error requesting access token from %s: %v", params.Party, err)
            }
            fmt.Printf("%s access token response: %v", params.Party, tokenResponse)
            if tokenResponse.AccessToken == "" {
                return "", errors.New("empty access token")
            }

            // Store the access token
            setTokenParams := SetTokenParams{
                TokenKind: "access",
                Party: params.Party,
                UserID: params.UserID, 
                Token: tokenResponse.AccessToken,
                ExpiresIn: tokenResponse.ExpiresIn,
                AppCtx: params.AppCtx,
            }
            if err := SetToken(setTokenParams); err != nil {
                return "", fmt.Errorf("error storing access token in Redis: %v", err)
            }

            // Store the new refresh token
            setTokenParams.TokenKind = "refresh"
            setTokenParams.Token = tokenResponse.RefreshToken
            setTokenParams.ExpiresIn = 0
            if err := SetToken(setTokenParams); err != nil {
                return "", fmt.Errorf("error storing refresh token in Redis: %v", err)
            }
            // Successfully requested a new access token from <party> using the refresh token
            return tokenResponse.AccessToken, nil
        } else {
            // Refresh token is expired, so user must reauthenticate
            if err := HandleLogout(params.Updater, clearTokenParams); err != nil {
                log.Printf("Error handling forced logout for user %s: %v", params.UserID, err)
            }
            return "", fmt.Errorf("reauthentication required with %s", params.Party)
        }
    }
}

// To mitigate circular dependency between auth0 and utils packages
type UserMetadataUpdater interface {
    UpdateUserMetadata(userID string, updatedAuthStatus map[string]interface{}) error
}

// Called when user clicks "Log Out of <party>" button on the user interface
func HandleLogout(updater UserMetadataUpdater, params ClearTokensParams) error {
    clearTokenParams := ClearTokensParams{ Party: params.Party, UserID: params.UserID, AppCtx: params.AppCtx}

    if err := ClearTokens(clearTokenParams); err != nil {
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