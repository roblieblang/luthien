package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

// TODO: Modularize this code a bit more. SRP. DRY (which you are, many times). Fix it.

type TokenResponse struct {
    AccessToken     string `json:"access_token"`
    TokenType       string `json:"token_type"`
    Scope           string `json:"scope"`
    ExpiresIn       int    `json:"expires_in"`
    RefreshToken    string `json:"refresh_token"`
}

// When hit, generates a code verifier and code challenge then redirects the user to Spotify auth URL
func LoginHandler(redisClient *redis.Client, c *gin.Context, clientID string, redirectURI string) {
    codeVerifier, err := utils.GenerateCodeVerifier(64)
    if err != nil {
        fmt.Printf("There was an issue generating the code verifier: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    err = redisClient.Set(c.Request.Context(), "codeVerifierKey", codeVerifier, time.Hour).Err()
    if err != nil {
        fmt.Printf("There was an issue storing the code verifier: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    codeChallenge := utils.SHA256Hash(codeVerifier)

    scope := "user-read-private user-read-email"
    params := url.Values{}
    params.Add("client_id", clientID)
    params.Add("response_type", "code")
    params.Add("redirect_uri", redirectURI)
    params.Add("scope", scope)
    params.Add("code_challenge_method", "S256")
    params.Add("code_challenge", codeChallenge)

    authURL := "https://accounts.spotify.com/authorize?" + params.Encode()

    c.JSON(http.StatusOK, gin.H{"authURL": authURL})

    fmt.Printf("Response sent: %v\n", gin.H{"authURL": authURL})
} 

func requestSpotifyToken(payload url.Values) (TokenResponse, error) {
    resp, err := http.PostForm("https://accounts.spotify.com/api/token", payload)
    if err != nil {
        return TokenResponse{}, err
    }
    defer resp.Body.Close()

    var tokenResponse TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return TokenResponse{}, err
    }

    return tokenResponse, nil
}


// Once user authorizes the application, Spotify redirects to a callback URL specified in application settings 
func CallbackHandler(redisClient *redis.Client, c *gin.Context, clientID string, redirectURI string) {
    // Extract auth code from query params
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
        return 
    }
    
    error := c.Query("error")
    if error != "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": error})
        return
    }

    codeVerifierCmd := redisClient.Get(c.Request.Context(), "codeVerifierKey")
    codeVerifier, err := codeVerifierCmd.Result()
    if err != nil {
        fmt.Printf("There was an issue retrieving the code verifier: %v\n", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    payload := url.Values{}
    payload.Set("grant_type", "authorization_code")
    payload.Set("code", code)
    payload.Set("redirect_uri", redirectURI) 
    payload.Set("client_id", clientID) 
    payload.Set("code_verifier", codeVerifier) 

    tokenResponse, err := requestSpotifyToken(payload)
    if err != nil {
        fmt.Printf("Error requesting access token from Spotify: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error requesting access token from Spotify"})
        return
    }

    // TODO: this is a temporarily hardcoded user ID
    userID := "65a88a1599c5cf91244826fb"

    // Store access and refresh tokens in Redis
    expiresIn := time.Duration(tokenResponse.ExpiresIn) * time.Second
    err = redisClient.Set(
        c.Request.Context(), 
        "accessToken:" + userID, 
        tokenResponse.AccessToken, 
        expiresIn).Err()
    if err != nil {
        fmt.Printf("There was an issue storing the access token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    expiresIn = time.Duration(tokenResponse.ExpiresIn) * time.Second
    err = redisClient.Set(
        c.Request.Context(), 
        "refreshToken:" + userID, 
        tokenResponse.RefreshToken, 
        expiresIn).Err()
    if err != nil {
        fmt.Printf("There was an issue storing the refresh token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    c.Redirect(http.StatusFound, "http://localhost:5173/")
}


func RefreshTokenHandler(redisClient *redis.Client, c *gin.Context, clientID string) {
    // Retrieve stored refreshToken
    // TODO: this is a temporarily hardcoded user ID
    userID := "65a88a1599c5cf91244826fb"
    refreshToken, err := redisClient.Get(c.Request.Context(), "refreshToken:" + userID).Result()
    if err == redis.Nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    payload := url.Values{}
    payload.Set("grant_type", "refresh_token")
    payload.Set("refresh_token", refreshToken)
    payload.Set("client_id", clientID)

    tokenResponse, err := requestSpotifyToken(payload)
    if err != nil {
        fmt.Printf("Error requesting refresh token from Spotify: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error requesting refresh token from Spotify"})
        return
    }

    expiresIn := time.Duration(tokenResponse.ExpiresIn) * time.Second
    err = redisClient.Set(
        c.Request.Context(), 
        "accessToken:" + userID, 
        tokenResponse.AccessToken, 
        expiresIn).Err()
    if err != nil {
        fmt.Printf("There was an issue storing the access token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    expiresIn = time.Duration(tokenResponse.ExpiresIn) * time.Second
    err = redisClient.Set(
        c.Request.Context(), 
        "refreshToken:" + userID, 
        tokenResponse.RefreshToken, 
        expiresIn).Err()
    if err != nil {
        fmt.Printf("There was an issue storing the refresh token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    

    c.JSON(http.StatusOK, gin.H{"access_token": tokenResponse.AccessToken})
}


func CheckAuthHandler(redisClient *redis.Client, c *gin.Context) {
    // TODO: Need a way to identify the current user
    // userID := getUserID(c)
    // TODO: this is a temporarily hardcoded user ID
    userID := "65a88a1599c5cf91244826fb"

     // Check if the access token exists for this user
     _, err := redisClient.Get(c.Request.Context(), "accessToken:" + userID).Result()
     if err == redis.Nil {
         // No access token found
         c.JSON(http.StatusOK, gin.H{"isAuthenticated": false})
         return
     } else if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
         return
     }
 
     // If an access token exists, the user is authenticated
     c.JSON(http.StatusOK, gin.H{"isAuthenticated": true})
}

func LogoutHandler(redisClient *redis.Client, c *gin.Context) {
    // TODO: this is a temporarily hardcoded user ID
    userID := "65a88a1599c5cf91244826fb"

    // Delete the access and refresh tokens from Redis
    _, err := redisClient.Del(c.Request.Context(), "accessToken:" + userID).Result()
    if err != nil {
        fmt.Printf("There was an issue deleting the access token: %v\n", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    _, err = redisClient.Del(c.Request.Context(), "refreshToken:" + userID).Result()
    if err != nil {
        fmt.Printf("There was an issue deleting the refresh token: %v\n", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}