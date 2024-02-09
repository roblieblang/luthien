package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/roblieblang/luthien/backend/internal/utils"
)

// Constructs and sends the direct HTTP requests to the Spotify API
// Handles OAuth authentication and token refresh
// Converts errors received by Spotify into a format usable by our app (if necessary)

type SpotifyClient struct {
    AppContext *utils.AppContext
}

type SpotifyTokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    Scope        string `json:"scope"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
}

func NewSpotifyClient(appCtx *utils.AppContext) *SpotifyClient {
    return &SpotifyClient{
        AppContext: appCtx,
    }
}

func (c *SpotifyClient) RequestToken(payload url.Values) (SpotifyTokenResponse, error) {
    resp, err := http.PostForm("https://accounts.spotify.com/api/token", payload)
    if err != nil {
        return SpotifyTokenResponse{}, err
    }
    defer resp.Body.Close()

    var tokenResponse SpotifyTokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return SpotifyTokenResponse{}, err
    }
    return tokenResponse, nil
}

// TODO: Get user profile

// TODO: Get user's playlists

// TODO: Get playlist

// TODO: Create playlist

// TODO: Update playlist 

// TODO: Delete playlist 


