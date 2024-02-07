package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

type Auth0TokenResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
}

func getAndSetManagementAPIAccessToken(c *gin.Context, appCtx *utils.AppContext) {
	clientSecret := appCtx.EnvConfig.Auth0ManagementClientSecret
	clientID := appCtx.EnvConfig.Auth0ManagementClientID
	domain := appCtx.EnvConfig.Auth0Domain

	url := "https://" + domain + "/oauth/token"

	payload := fmt.Sprintf(
		"{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"https://%s/api/v2/\",\"grant_type\":\"client_credentials\"}", 
		clientID, 
		clientSecret, 
		domain,
	)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
        log.Printf("There was an issue requesting the Auth0 Management API access token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
        log.Printf("There was an issue reading the response body: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

	// Unmarshal the JSON response into the Auth0TokenResponse struct
    var tokenResponse Auth0TokenResponse
    err = json.Unmarshal(body, &tokenResponse)
    if err != nil {
        log.Printf("There was an issue unmarshaling the response: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

	// Store access token in Redis
	err = appCtx.RedisClient.Set(
		c.Request.Context(), 
		"auth0ManagementAPIAccessToken", 
		tokenResponse.AccessToken, 
		time.Duration(tokenResponse.ExpiresIn) * time.Second).Err()
    if err != nil {
        log.Printf("There was an issue storing the Auth0 Management API Access Token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
}

// Helper function that orchestrates all work around getting a usable access token
func getAccessTokenHelper(c *gin.Context, appCtx *utils.AppContext) string {
	tokenExpired, err := utils.IsAccessTokenExpired(c, appCtx, "auth0ManagementAPIAccessToken")
	if err != nil {
		log.Printf("There was an error checking token expiry: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
        return ""
	}

	if tokenExpired {
		getAndSetManagementAPIAccessToken(c, appCtx)
	} 
	
	// Retrieve access token from Redis
	accessToken, err := appCtx.RedisClient.Get(c.Request.Context(), "auth0ManagementAPIAccessToken").Result()
	if err != nil {
		log.Printf("Failed to retrieve Auth0 Management API Access Token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return ""
	}
	return accessToken
}

// Updates user's app metadata to indicate whether they have authenticated with Spotify
func UpdateSpotifyAuthStatus(c *gin.Context, appCtx *utils.AppContext, userID string, isAuthenticated bool) {
	accessToken := getAccessTokenHelper(c, appCtx)

	domain := appCtx.EnvConfig.Auth0Domain
	url := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userID)
	
	metadata := map[string]interface{}{
		"app_metadata": map[string]bool{
			"authenticated_with_spotify": isAuthenticated,
		},
	}
	payload, err := json.Marshal(metadata)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	if err != nil {
        // c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HTTP request"})
		log.Printf("Failed to create HTTP request: %v", err)
        return
    }

	httpReq.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))
    httpReq.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(httpReq)
    if err != nil {
		log.Printf("Failed to execute request: %v", err)
        return
    }
    defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
    if err != nil {
		log.Printf("Failed to read response body: %v", err)
        return
    }

	if res.StatusCode >= 400 {
		log.Printf("Received error status from Auth0: %s", string(body))
        return
    }

	log.Printf("Spotify authentication status updated successfully for user %s", userID)
}

type Auth0UserMetadata struct {
	CreatedAt   time.Time `json:"created_at"`
	Email       string    `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	Identities  []Identity `json:"identities"`
	Name        string    `json:"name"`
	Nickname    string    `json:"nickname"`
	Picture     string    `json:"picture"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id"`
	AppMetadata struct {
		AuthenticatedWithSpotify bool `json:"authenticated_with_spotify"`
	} `json:"app_metadata"`
	LastIP      string    `json:"last_ip"`
	LastLogin   time.Time `json:"last_login"`
	LoginsCount int       `json:"logins_count"`
}
type Identity struct {
	Connection string `json:"connection"`
	Provider   string `json:"provider"`
	UserID     string `json:"user_id"`
	IsSocial   bool   `json:"isSocial"`
}

// Retrieves a user's Auth0 metadata
func GetUserMetadata(c *gin.Context, appCtx *utils.AppContext, userID string) (Auth0UserMetadata, error){
	domain := appCtx.EnvConfig.Auth0Domain
	url := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
        return Auth0UserMetadata{}, err
    }

	accessToken := getAccessTokenHelper(c, appCtx)

	httpReq.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))
    httpReq.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		fmt.Println(err)
		return Auth0UserMetadata{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
    if err != nil {
		log.Printf("Failed to read response body: %v", err)
        return Auth0UserMetadata{}, err
    }

	if res.StatusCode >= 400 {
		log.Printf("Received error status from Auth0: %s", string(body))
        return Auth0UserMetadata{}, err
    }

	var userMetadata Auth0UserMetadata

	if err := json.Unmarshal(body, &userMetadata); err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		return Auth0UserMetadata{}, err
	}

	return userMetadata, nil
}


// TODO: I forgot what this is for...
func ValidateAuth0LoginToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: extract and validate login token. Is it actually called a "login token"?
		
	})
}
