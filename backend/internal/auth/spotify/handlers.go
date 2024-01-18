package spotify

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien-core-server/internal/utils"
)

// When hit, generates a code verifier and code challenge then redirects the user to Spotify auth URL
func LoginHandler(c *gin.Context, clientID string, redirectURI string) {
    // TODO: store code verifier in a session or some other secure place (REDIS!!!) redis redis (REDIS) redis
    codeVerifier, err := utils.GenerateCodeVerifier(64)
    if err != nil {
        fmt.Printf("There was an issue generating the code verifier: %v", err)
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

// Once user authorizes the application, Spotify redirects to a callback URL specified in application settings 
func CallbackHandler(c *gin.Context, clientID string, redirectURI string) {
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

    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Set("code", code)
    data.Set("redirect_uri", redirectURI) 
    data.Set("client_id", clientID) 
    data.Set("code_verifier", ) 

    resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request Spotify access token"})
        return
    }
    defer resp.Body.Close()
    // TODO: store code verifier in a session with Redis and retrieve it from there
    //TODO: first dockerize the whole app so that you can run Redis outside of WSL (less config for all devs involved)


    }

// 3. get refresh tokens
func RefreshTokenHandler(c *gin.Context) {
    // Retrieve stored refreshToken
    // accessToken, newRefreshToken := refreshSpotifyToken(refreshToken)

    // // Update stored refreshToken and send new accessToken to the frontend
    // json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
}