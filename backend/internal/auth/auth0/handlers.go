package auth0

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/roblieblang/luthien/backend/internal/utils"
	// "github.com/golang-jwt/jwt/v5"
)

type Auth0TokenResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
}

func GetManagementAPIAcessToken(redisClient *redis.Client, c *gin.Context) {
	envConfig := utils.LoadENV()

	clientId := envConfig.Auth0ManagementClientID
	clientSecret := envConfig.Auth0ManagementClientSecret
	auth0Domain := envConfig.Auth0Domain

	url := "https://" + auth0Domain + "/oauth/token"

	payload := fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"https://%s/api/v2/\",\"grant_type\":\"client_credentials\"}", clientId, clientSecret, auth0Domain)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
        fmt.Printf("There was an issue requesting the Auth0 Management API access token: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
        fmt.Printf("There was an issue reading the response body: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

	// Unmarshal the JSON response into the Auth0TokenResponse struct
    var tokenResponse Auth0TokenResponse
    err = json.Unmarshal(body, &tokenResponse)
    if err != nil {
        fmt.Printf("There was an issue unmarshaling the response: %v", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

	fmt.Printf("Access Token: %s, Expires In: %d\n", tokenResponse.AccessToken, tokenResponse.ExpiresIn)
	// TODO: store in Redis and check expiration time units


}

func ValidateAuth0LoginToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: extract and validate login token. Is it actually called a "login token"?
		
	})
}


// func UpdateAuth0UserMetadata(auth0UserID string, metadata interface{}) error {
// 	
// }