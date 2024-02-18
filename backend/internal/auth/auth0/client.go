package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/roblieblang/luthien/backend/internal/utils"
)

// Auth0Client manages communication with the Auth0 Management API.
type Auth0Client struct {
    AppContext *utils.AppContext
}

type Auth0TokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
	Scope 		string `json:"scope"`
	TokenType 	string `json:"token_type"`
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
	Connection 	string `json:"connection"`
	Provider   	string `json:"provider"`
	AccessToken string `json:"access_token"`
	ExpiresIn	int	   `json:"expires_in"`
	UserID     	string `json:"user_id"`
	IsSocial   	bool   `json:"isSocial"`
}

func NewAuth0Client(appCtx *utils.AppContext) *Auth0Client {
    return &Auth0Client{
        AppContext: appCtx,
    }
}

// Requests a new Auth0 Management API access token
func (c *Auth0Client) RequestToken() (Auth0TokenResponse, error) {
    clientSecret := c.AppContext.EnvConfig.Auth0ManagementClientSecret
	clientID := c.AppContext.EnvConfig.Auth0ManagementClientID
	domain := c.AppContext.EnvConfig.Auth0Domain

    url_ := "https://"+domain+"/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("audience", "https://" + domain + "/api/v2/")

	payload := strings.NewReader(data.Encode())


    req, err := http.NewRequest("POST", url_, payload)
    if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
        return Auth0TokenResponse{}, err
    }

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

    res, err := http.DefaultClient.Do(req)
	if err != nil {
        log.Printf("There was an issue requesting the Auth0 Management API access token: %v", err)
        return Auth0TokenResponse{}, err
    }
	defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
	if err != nil {
        log.Printf("There was an issue reading the response body: %v", err)
        return Auth0TokenResponse{}, err
    }

    // Unmarshal the JSON response into the Auth0TokenResponse struct
    var tokenResponse Auth0TokenResponse
    err = json.Unmarshal(body, &tokenResponse)
    if err != nil {
        log.Printf("There was an issue unmarshaling the response: %v", err)
        return Auth0TokenResponse{}, err
    }

    return tokenResponse, nil
}

func (c *Auth0Client) GetUserMetadata(accessToken string, userID string) (Auth0UserMetadata, error) {
    domain := c.AppContext.EnvConfig.Auth0Domain
	url := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userID)
    
    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
        return Auth0UserMetadata{}, err
    }

	req.Header.Add("Accept", "application/json")
    req.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
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

// Make a partial update of a user's metadata
func (c *Auth0Client) UpdateUserMetadata(accessToken, userID string, metadata map[string]interface{}) error {
    domain := c.AppContext.EnvConfig.Auth0Domain
	url := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userID)

    updatedFields := make([]string, 0, len(metadata))
    for field := range metadata {
        updatedFields = append(updatedFields, field)
    }
    log.Printf("Updating metadata fields %v for user %s", updatedFields, userID)

    payload, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return err
	}

    req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
        return err
    }

    req.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Add("content-type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
		log.Printf("Failed to execute request: %v", err)
        return err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
		log.Printf("Failed to read response body: %v", err)
        return err
    }

    if res.StatusCode >= 400 {
		log.Printf("Received error status from Auth0: %s", string(body))
        return err
    }

    log.Printf("Metadata updated successfully for user %s: %v", userID, updatedFields)
	return nil
}