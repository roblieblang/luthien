package spotify

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

// TODO: Modularize this code a bit more. SRP. DRY (which you are, many times). Fix it.
// TODO: reorganize this file into a services.go, ...,
type TokenResponse struct {
    AccessToken     string `json:"access_token"`
    TokenType       string `json:"token_type"`
    Scope           string `json:"scope"`
    ExpiresIn       int    `json:"expires_in"`
    RefreshToken    string `json:"refresh_token"`
}

// When hit, generates a code verifier and code challenge then redirects the user to Spotify auth URL
func LoginHandler(c *gin.Context, appCtx *utils.AppContext) {
    sessionID := utils.GenerateSessionID()
    codeVerifier, err := utils.GenerateCodeVerifier(64)
    if err != nil {
        log.Printf("There was an issue generating the code verifier: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    err = appCtx.RedisClient.Set(c.Request.Context(), "spotifyCodeVerifier:" + sessionID, codeVerifier, time.Minute * 10).Err()
    if err != nil {
        log.Printf("There was an issue storing the code verifier: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    codeChallenge := utils.SHA256Hash(codeVerifier)

    scope := "user-read-private user-read-email"
    params := url.Values{}
    params.Add("client_id", appCtx.EnvConfig.SpotifyClientID)
    params.Add("response_type", "code")
    params.Add("redirect_uri", appCtx.EnvConfig.SpotifyRedirectURI)
    params.Add("scope", scope)
    params.Add("code_challenge_method", "S256")
    params.Add("code_challenge", codeChallenge)

    authURL := "https://accounts.spotify.com/authorize?" + params.Encode()

    c.JSON(http.StatusOK, gin.H{"authURL": authURL, "sessionID": sessionID})

    log.Printf("Response sent: %v\n", gin.H{"authURL": authURL, "sessionID": sessionID})
} 

// Requests a new access token and refresh token from Spotify
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
// This is called by Spotify, not our own application 
// Part of PKCE flow
func CallbackHandler(c *gin.Context, appCtx *utils.AppContext) {
    var req struct {
        Code   string `json:"code"`
        UserID string `json:"userID"`
        SessionID string `json:"sessionID"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    if req.Code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
        return
    }
    if req.UserID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }
    if req.SessionID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
        return
    }

    codeVerifier, err := appCtx.RedisClient.Get(c.Request.Context(), "spotifyCodeVerifier:" + req.SessionID).Result()
    if err != nil {
        log.Printf("There was an issue retrieving the code verifier: %v\n", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    payload := url.Values{}
    payload.Set("grant_type", "authorization_code")
    payload.Set("code", req.Code)
    payload.Set("redirect_uri", appCtx.EnvConfig.SpotifyRedirectURI) 
    payload.Set("client_id", appCtx.EnvConfig.SpotifyClientID) 
    payload.Set("code_verifier", codeVerifier) 

    tokenResponse, err := requestSpotifyToken(payload)
    if err != nil {
        log.Printf("Error requesting access token from Spotify: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error requesting access token from Spotify"})
        return
    }

    expiresIn := time.Duration(tokenResponse.ExpiresIn) * time.Second
    if tokenResponse.AccessToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Empty access token"})
        return
    }

    // Store access and refresh tokens in Redis
    err = appCtx.RedisClient.Set(
        c.Request.Context(), 
        "spotifyAccessToken:" + req.UserID, 
        tokenResponse.AccessToken, 
        expiresIn).Err()
    if err != nil {
        log.Printf("There was an issue storing the access token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    err = appCtx.RedisClient.Set(
        c.Request.Context(), 
        "spotifyRefreshToken:" + req.UserID, 
        tokenResponse.RefreshToken, 
        expiresIn).Err()
    if err != nil {
        log.Printf("There was an issue storing the refresh token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    // Update Auth0 user app_metadata 
    auth0.UpdateSpotifyAuthStatus(c, appCtx, req.UserID, true)
    c.JSON(http.StatusOK, gin.H{"redirectURL": "http://localhost:5173/"})
}


func refreshAccessToken(c *gin.Context, appCtx *utils.AppContext, userID string) {
    // Retrieve stored refresh token
    refreshToken, err := appCtx.RedisClient.Get(c.Request.Context(), "spotifyRefreshToken:" + userID).Result()
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
    payload.Set("client_id", appCtx.EnvConfig.SpotifyClientID)

    tokenResponse, err := requestSpotifyToken(payload)
    if err != nil {
        log.Printf("Error requesting refresh token from Spotify: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error requesting refresh token from Spotify"})
        return
    }

    expiresIn := time.Duration(tokenResponse.ExpiresIn) * time.Second
    err = appCtx.RedisClient.Set(
        c.Request.Context(), 
        "spotifyAccessToken:" + userID, 
        tokenResponse.AccessToken, 
        expiresIn).Err()
    if err != nil {
        log.Printf("There was an issue storing the access token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    expiresIn = time.Duration(tokenResponse.ExpiresIn) * time.Second
    err = appCtx.RedisClient.Set(
        c.Request.Context(), 
        "spotifyRefreshToken:" + userID, 
        tokenResponse.RefreshToken, 
        expiresIn).Err()
    if err != nil {
        log.Printf("There was an issue storing the refresh token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    

    c.JSON(http.StatusOK, gin.H{"access_token": tokenResponse.AccessToken})
}

// Checks Spotify authentication status for a specific user
func CheckAuthHandler(c *gin.Context, appCtx *utils.AppContext) {
    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    userMetadata, err := auth0.GetUserMetadata(c, appCtx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user metadata"})
		return
    }

    c.JSON(http.StatusOK, gin.H{"isAuthenticated": userMetadata.AppMetadata.AuthenticatedWithSpotify})
}

func LogoutHandler(c *gin.Context, appCtx *utils.AppContext) {
    var req struct {
        UserID string `json:"userID"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
    }  
    
    userID := req.UserID
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    // Delete the access and refresh tokens from Redis
    _, err := appCtx.RedisClient.Del(c.Request.Context(), "spotifyAccessToken:" + userID).Result()
    if err != nil {
        log.Printf("There was an issue deleting the access token: %v\n", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    _, err = appCtx.RedisClient.Del(c.Request.Context(), "spotifyRefreshToken:" + userID).Result()
    if err != nil {
        log.Printf("There was an issue deleting the refresh token: %v\n", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    
    // Update Auth0 user app_metadata 
    auth0.UpdateSpotifyAuthStatus(c, appCtx, userID, false)

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}