package spotify

import (
	"bytes"
	"encoding/json"
    "encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
    "strings"
	"net/url"

	"github.com/roblieblang/luthien/backend/internal/utils"
)

// Constructs and sends the direct HTTP requests to the Spotify API
// Handles OAuth authentication and token refresh
// Converts errors received by Spotify into a format usable by our app (if necessary)

type SpotifyClient struct {
    AppContext *utils.AppContext
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
    // Height int    `json:"height"` // don't need these fields at present
    // Width  int    `json:"width"`
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
    Images        []Image          `json:"images"`  // may be empty or contain up to three images
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

type SpotifyPlaylistTracksResponse struct {
    Limit  int                 `json:"limit"`
    Offset int                 `json:"offset"`
    Items  []PlaylistTrackItem `json:"items"`
}

type PlaylistTrackItem struct {
    Track TrackDetails `json:"track"`
}

type TrackDetails struct {
    Album       AlbumDetails   `json:"album"`
    Artists     []Artist       `json:"artists"`
    ExternalIDs ExternalIDs    `json:"external_ids"`
    Name        string         `json:"name"`
}

type AlbumDetails struct {
    Name        string  `json:"name"`
    Images      []Image `json:"images"`
}

type Artist struct {
    Name string `json:"name"`
}

type ExternalIDs struct {
    ISRC string `json:"isrc"`
}

type CreatePlaylistBody struct {
    UserID        string                `json:"userId"`
    SpotifyUserID string                `json:"spotifyUserId"`
    Payload       CreatePlaylistPayload `json:"payload"`
}

type CreatePlaylistPayload struct {
    Name          string `json:"name"`
    Public        *bool  `json:"public,omitempty"`  // Spotify's API defaults this to true
    Collaborative *bool  `json:"collaborative,omitempty"`  // Defaults to false, to be true public must be false
    Description   string `json:"description,omitempty"`
}

type SpotifySearchResponse struct {
    Tracks struct {
        Items []struct {
            URI string `json:"uri"`
        } `json:"items"`
    } `json:"tracks"`
}

// Returns a new SpotifyClient struct 
func NewSpotifyClient(appCtx *utils.AppContext) *SpotifyClient {
    return &SpotifyClient{
        AppContext: appCtx,
    }
}

// Requests a new access token from Spotify
func (c *SpotifyClient) RequestToken(payload url.Values) (utils.TokenResponse, error) {
    client := &http.Client{}
    req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
    if err != nil {
        log.Printf("Error creating request for Spotify token: %v\n", err)
        return utils.TokenResponse{}, err
    }

    // Add Authorization header if this is a refresh token request
    if payload.Get("grant_type") == "refresh_token" {
        authHeaderVal := base64.StdEncoding.EncodeToString([]byte(c.AppContext.EnvConfig.SpotifyClientID + ":" + c.AppContext.EnvConfig.SpotifyClientSecret))
        req.Header.Add("Authorization", "Basic "+authHeaderVal)
    }

    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(req)

    if resp.StatusCode != http.StatusOK {
        var errorResponse struct {
            Error            string `json:"error"`
            ErrorDescription string `json:"error_description"`
        }
        if decodeErr := json.NewDecoder(resp.Body).Decode(&errorResponse); decodeErr == nil {
            log.Printf("Spotify API error: %s - %s\n", errorResponse.Error, errorResponse.ErrorDescription)
        } else {
            log.Printf("Failed to decode Spotify error response: %v\n", decodeErr)
        }
        return utils.TokenResponse{}, fmt.Errorf("spotify API request failed: %s", resp.Status)
    }

    if err != nil {
        log.Printf("Error requesting token from Spotify: %v\n", err)
        return utils.TokenResponse{}, err
    }
    defer resp.Body.Close()
    
    var tokenResponse utils.TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        log.Printf("Error reading Spotify token response: %v\n", err)
        return utils.TokenResponse{}, err
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

// Builds a URL with necessary fields and params for Spotify's Get Playlist Items endpoint
func (c *SpotifyClient) buildPlaylistItemsURL(playlistID string, limit, offset int) string {
    baseURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

    // If you end up needing the fields param:
    // items(track(name,external_ids.isrc,artists(name),album(name,images))),limit,offset

    params := url.Values{}
    params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))

    fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

    return fullURL
}

// Gets playlist items (with pagination [necessary due to API rate limits]).
// Spotify documentation: https://developer.spotify.com/documentation/web-api/reference/get-playlists-tracks
func (c *SpotifyClient) GetPlaylistTracks(accessToken, playlistID string, limit, offset int) (SpotifyPlaylistTracksResponse, error) {
    url := c.buildPlaylistItemsURL(playlistID, limit, offset)

    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
        return SpotifyPlaylistTracksResponse{}, err
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SpotifyPlaylistTracksResponse{}, err
	}
	defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return SpotifyPlaylistTracksResponse{}, err
    }

    if res.StatusCode >= 400 {
        return SpotifyPlaylistTracksResponse{}, fmt.Errorf("spotify API request failed with status %d: %s", res.StatusCode, string(body))
    }

	var playlistTracks SpotifyPlaylistTracksResponse
	if err := json.Unmarshal(body, &playlistTracks); err != nil {
		return SpotifyPlaylistTracksResponse{}, err
	}

    return playlistTracks, nil
}

// Creates a new playlist
func (c *SpotifyClient) CreatePlaylist(accessToken, spotifyUserID string, playlistPayload CreatePlaylistPayload) ([]byte, error) {
    url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", spotifyUserID)
    
    payload, err := json.Marshal(playlistPayload)
    if err != nil {
        // TODO: standardize error handling across the application
        return nil, fmt.Errorf("error marshaling payload: %w", err)
	}
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Add("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    if res.StatusCode >= 400 {
        return nil, fmt.Errorf("spotify API error (status code %d): %s", res.StatusCode, string(body))
    }
    
    return body, nil
}

/* YOUTUBE TO SPOTIFY CONVERSION FLOW */
// 1. Get ISRC of all items in YouTube playlist (or similar identifier, maybe just combination of artist, album, title)
//      - Can get ISRC with MusicBrainz API by searching for a recording with a title that matches the YouTube video title
//      - https://musicbrainz.org/doc/MusicBrainz_API
//      - This might be overkill, so maybe just do a simple search with cleaned video titles and see how that goes
// 2. Use ISRC of each playlist item to find the equivalent item on Spotify and get its Spotify URI
//    using GetURIWithISRC
// 3. Assemble all URIs for the playlist-in-conversion into an array of strings
// 4. Create a new Spotify playlist for the conversion
// 5. Pass the URI array into the body of an AddItemsToPlaylist request

// TODO: implement
// Docs: https://developer.spotify.com/documentation/web-api/reference/search
// Searches for a specific Spotify track using an ISRC 
// (will we need a fallback? e.g. combination of artist, album, track, year)
func (c *SpotifyClient) GetURIWithISRC(accessToken, isrc string) (SpotifySearchResponse, error) {    
    return SpotifySearchResponse{}, nil
}

// TODO: implement 
// Docs: https://developer.spotify.com/documentation/web-api/reference/add-tracks-to-playlist
// Adds items to an existing Spotify playlist
func (c *SpotifyClient) AddItemsToPlaylist(accessToken, playlistID string) error {
    return nil
}

// TODO: Implement Update Playlist(?) 

