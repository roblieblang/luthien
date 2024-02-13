package spotify

import (
	"encoding/json"
	"fmt"
	"io"
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

type SpotifyUserProfile struct {
    Country        string            `json:"country"`
    DisplayName    string            `json:"display_name"`
    Email          string            `json:"email"`
    ExplicitContent ExplicitContent  `json:"explicit_content"`
    ExternalUrls   ExternalUrls      `json:"external_urls"`
    Followers      Followers         `json:"followers"`
    Href           string            `json:"href"`
    ID             string            `json:"id"`
    Images         []Image           `json:"images"`
    Product        string            `json:"product"`
    Type           string            `json:"type"`
    URI            string            `json:"uri"`
}

type ExplicitContent struct {
    FilterEnabled bool `json:"filter_enabled"`
    FilterLocked  bool `json:"filter_locked"`
}

type ExternalUrls struct {
    Spotify string `json:"spotify"`
}

type Followers struct {
    Href  string `json:"href"`
    Total int    `json:"total"`
}

type Image struct {
    URL    string `json:"url"`
    Height int    `json:"height"`
    Width  int    `json:"width"`
}

type SpotifyPlaylistsResponse struct {
    Href     string          `json:"href"`  // link to endpoint returning full result of the request
    Items    []PlaylistItem  `json:"items"`
    Limit    int             `json:"limit"`  // max number of itesm in the response
    Next     *string          `json:"next"`  // URL to the next page of items (null if none)
    Offset   int             `json:"offset"`
    Previous *string         `json:"previous"`  // URL to the previous page of items (null if none)
    Total    int             `json:"total"`  // total number of items (playlists) for the user
}

type PlaylistItem struct {
    ID            string           `json:"id"`
    Name          string           `json:"name"`
    Description   string           `json:"description"`
    Collaborative bool             `json:"collaborative"`  // true if owner allows other users to modify the playlist
    Images        []PlaylistImage  `json:"images"`  // may be empty or contain up to three images
    ExternalUrls  ExternalUrls     `json:"external_urls"`
    Href          string           `json:"href"`
    Owner         PlaylistOwner    `json:"owner"`
    PrimaryColor  *string          `json:"primary_color"` // Using *string to allow for null value
    Public        bool             `json:"public"`
    SnapshotID    string           `json:"snapshot_id"`
    Tracks        PlaylistTracks   `json:"tracks"`  // collection containing a link to endpoint which contains full details of playlist tracks
    Type          string           `json:"type"`
    URI           string           `json:"uri"`
}

type PlaylistImage struct {
    Height *int    `json:"height"` // Using *int to allow for null value
    URL    string  `json:"url"`
    Width  *int    `json:"width"` // Using *int to allow for null value
}

type PlaylistOwner struct {
    DisplayName  string        `json:"display_name"`
    ExternalUrls ExternalUrls  `json:"external_urls"`
    Href         string        `json:"href"`
    ID           string        `json:"id"`
    Type         string        `json:"type"`
    URI          string        `json:"uri"`
}

type PlaylistTracks struct {
    Href  string `json:"href"`  // link to endpoint which contains full details for all tracks in the playlist
    Total int    `json:"total"`
}

func NewSpotifyClient(appCtx *utils.AppContext) *SpotifyClient {
    return &SpotifyClient{
        AppContext: appCtx,
    }
}

// Requests a new access token from Spotify
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

// Gets the current user's profile
func (c *SpotifyClient) GetCurrentUserProfile(accessToken string) (SpotifyUserProfile, error) {
    url := "https://api.spotify.com/v1/me"

    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
        return SpotifyUserProfile{}, err
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SpotifyUserProfile{}, err
	}
	defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return SpotifyUserProfile{}, err
    }

    if res.StatusCode >= 400 {
        return SpotifyUserProfile{}, fmt.Errorf("spotify API request failed with status %d: %s", res.StatusCode, string(body))
    }

	var userProfile SpotifyUserProfile
	if err := json.Unmarshal(body, &userProfile); err != nil {
		return SpotifyUserProfile{}, err
	}

	return userProfile, nil
}

// Gets the current user's playlists
func (c *SpotifyClient) GetCurrentUserPlaylists(accessToken string, limit, offset int) (SpotifyPlaylistsResponse, error) {
    url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", limit, offset)

    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
        return SpotifyPlaylistsResponse{}, err
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SpotifyPlaylistsResponse{}, err
	}
	defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return SpotifyPlaylistsResponse{}, err
    }

    if res.StatusCode >= 400 {
        return SpotifyPlaylistsResponse{}, fmt.Errorf("spotify API request failed with status %d: %s", res.StatusCode, string(body))
    }

	var playlistsResponse SpotifyPlaylistsResponse
	if err := json.Unmarshal(body, &playlistsResponse); err != nil {
		return SpotifyPlaylistsResponse{}, err
	}

	return playlistsResponse, nil
}

// TODO: Gets playlist tracks https://developer.spotify.com/documentation/web-api/reference/get-playlist
func (c *SpotifyClient) GetPlaylistTracks(accessToken, playlistID string, limit, offset int) (SpotifyPlaylistsResponse, error) {
    return SpotifyPlaylistsResponse{}, nil
}

// TODO: Implement Create playlist

// TODO: Implement Update playlist 

