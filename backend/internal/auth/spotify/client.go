package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

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
    Limit    int             `json:"limit"`  // max number of itesm in the response
    Next     *string          `json:"next"`  // URL to the next page of items (null if none)
    Offset   int             `json:"offset"`
    Previous *string         `json:"previous"`  // URL to the previous page of items (null if none)
    Total    int             `json:"total"`  // total number of items (playlists) for the user
    Items    []PlaylistItem  `json:"items"`
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
    URI         string         `json:"uri"`
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
func (c *SpotifyClient) GetCurrentUserPlaylists(accessToken string, offset int) (SpotifyPlaylistsResponse, error) {
    url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", 20, offset) 

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

// Gets playlist items
func (c *SpotifyClient) GetPlaylistTracks(accessToken, playlistID string) (SpotifyPlaylistTracksResponse, error) {
    var allTracks []PlaylistTrackItem
    var limit = 100
    var offset = 0

    for {
        url := c.buildPlaylistItemsURL(playlistID, limit, offset)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating new request: %v", err)
            return SpotifyPlaylistTracksResponse{}, err
        }
        req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            log.Printf("Error making request to Spotify: %v", err)
            return SpotifyPlaylistTracksResponse{}, err
        }
        defer res.Body.Close()

        if res.StatusCode >= 400 {
            body, err := io.ReadAll(res.Body) 
            if err != nil {
                return SpotifyPlaylistTracksResponse{}, fmt.Errorf("error reading response body: %w", err)
            }
            return SpotifyPlaylistTracksResponse{}, fmt.Errorf("Spotify API request failed with status %d: %s", res.StatusCode, string(body))
        }

        var page SpotifyPlaylistTracksResponse
        if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
            return SpotifyPlaylistTracksResponse{}, err
        }

        allTracks = append(allTracks, page.Items...)

        if len(page.Items) < limit {
            // Break the loop if this page has fewer items than the max, indicating it's the last page
            break
        }
        offset += limit // Prepare the offset for the next page
    }

    return SpotifyPlaylistTracksResponse{
        Items: allTracks,
    }, nil
}


// Creates a new playlist
func (c *SpotifyClient) CreatePlaylist(accessToken, spotifyUserID string, playlistPayload CreatePlaylistPayload) (string, error) {
    url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", spotifyUserID)
    
    payload, err := json.Marshal(playlistPayload)
    if err != nil {
        return "", fmt.Errorf("error marshaling payload: %w", err)
	}
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Add("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("error executing request: %w", err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    if res.StatusCode >= 400 {
        var spotifyError struct {
            Error struct {
                Status  int    `json:"status"`
                Message string `json:"message"`
            } `json:"error"`
        }
        json.Unmarshal(body, &spotifyError) 
        log.Printf("spotify API error: %d %s", spotifyError.Error.Status, spotifyError.Error.Message)
        return "", fmt.Errorf("spotify API error: %d %s", spotifyError.Error.Status, spotifyError.Error.Message)
    }
    var responseBody struct {
        ID string `json:"id"`
    }
    
    err = json.Unmarshal(body, &responseBody)
    if err != nil {
        return "", fmt.Errorf("error unmarshaling response body: %w", err)
    }

    return responseBody.ID, nil
}

type AddItemsToPlaylistPayload struct {
    ItemURIs []string   `json:"uris"`
    Position int        `json:"position"`
}

// Adds items to an existing Spotify playlist
func (c *SpotifyClient) AddItemsToPlaylist(accessToken, playlistID string, addItemsPayload AddItemsToPlaylistPayload) error {
    const maxItemsPerRequest = 100

    // Split ItemURIs into chunks of up to 100
    for i := 0; i < len(addItemsPayload.ItemURIs); i += maxItemsPerRequest {
        end := i + maxItemsPerRequest
        if end > len(addItemsPayload.ItemURIs) {
            end = len(addItemsPayload.ItemURIs)
        }
        
        chunk := AddItemsToPlaylistPayload{
            ItemURIs: addItemsPayload.ItemURIs[i:end],
            Position: addItemsPayload.Position,
        }
        
        if err := c.addItemsChunkToPlaylist(accessToken, playlistID, chunk); err != nil {
            return err
        }
    }

    return nil
}

// Helper function to add a chunk of items to the playlist
func (c *SpotifyClient) addItemsChunkToPlaylist(accessToken, playlistID string, payload AddItemsToPlaylistPayload) error {
    url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("error marshaling payload: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return fmt.Errorf("error creating request: %w", err)
    }
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Add("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("error executing request: %w", err)
    }
    defer res.Body.Close()

    if res.StatusCode >= 400 {
        body, _ := io.ReadAll(res.Body)
        return fmt.Errorf("spotify API error (status code %d): %s", res.StatusCode, string(body))
    }

    return nil
}

type SpotifySearchResponse struct {
    Tracks struct {
        Items []struct {
            Album struct {
                Name   string `json:"name"`
                Images []Image `json:"images"`
            } `json:"album"`
            Artists []struct {
                Name string `json:"name"`
            } `json:"artists"`
            Name string `json:"name"`
            URI  string `json:"uri"`
        } `json:"items"`
    } `json:"tracks"`
}

type Track struct {
    URI            string `json:"uri"`
    Title          string `json:"title"`
    ArtistNames    string `json:"artistNames"` // Concatenated names of the artists
    AlbumName      string `json:"albumName"`
    AlbumImageUrl  string `json:"albumImageUrl"`
}

// Processes a Spotify track search response and returns a struct containing relevant track information
func processSpotifySearchResponse(response SpotifySearchResponse) []utils.UnifiedTrackSearchResult {
    var searchResults []utils.UnifiedTrackSearchResult
    for _, item := range response.Tracks.Items {
        var albumImageURL string
        if len(item.Album.Images) > 0 {
            albumImageURL = item.Album.Images[0].URL
        }

        artistNames := make([]string, len(item.Artists))
        for i, artist := range item.Artists {
            artistNames[i] = artist.Name
        }
        searchResult := utils.UnifiedTrackSearchResult{
            ID:         item.URI,
            Title:      item.Name,
            Artist:     strings.Join(artistNames, ", "),
            Album:      item.Album.Name,
            Thumbnail:  albumImageURL,
        }

        searchResults = append(searchResults, searchResult)
    }
    return searchResults
}

// Builds a search URL to be used to search for matching Spotify tracks
func (c *SpotifyClient) buildSearchURL(artistName, trackTitle, videoTitle string, limit, offset int) string {
    baseURL := "https://api.spotify.com/v1/search"
    var queryParts []string

    if videoTitle != "" {
        encodedVideoTitle := url.QueryEscape(videoTitle)
        queryParts = append(queryParts, encodedVideoTitle)
    }

    if trackTitle != "" {
        encodedTrack := url.QueryEscape(trackTitle)
        queryParts = append(queryParts, fmt.Sprintf("track:%s", encodedTrack))
    }

    if artistName != "" {
        encodedArtist := url.QueryEscape(artistName)
        queryParts = append(queryParts, fmt.Sprintf("artist:%s", encodedArtist))
    }

    query := strings.Join(queryParts, " ")

    params := url.Values{}
    params.Add("query", query)
    params.Add("type", "track")
    params.Add("limit", fmt.Sprintf("%d", limit))
    params.Add("offset", fmt.Sprintf("%d", offset))

    fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
    return fullURL
}


// Retrieves Spotify tracks that match the provided artist name and track title
func (c *SpotifyClient) SearchTracksUsingArtistAndTrack(accessToken, artistName, trackTitle string, limit, offset int) ([]utils.UnifiedTrackSearchResult, error) {    
    url := c.buildSearchURL(artistName, trackTitle, "",  limit, offset)    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request %w", err)
    }

    log.Printf("Sending request to Spotify API: %s", url)

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }

    log.Printf("Received Spotify API response with status code: %d", res.StatusCode)
    if res.StatusCode != http.StatusOK {
        bodyBytes, err := io.ReadAll(res.Body)
        if err == nil {
            log.Printf("Spotify API response body: %s", string(bodyBytes))
        }
    }
    
    defer res.Body.Close()

    var response SpotifySearchResponse
    if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("error decoding response: %w", err)
    }

    if len(response.Tracks.Items) == 0 {
        return nil, fmt.Errorf("no tracks found for '%s' by '%s'", trackTitle, artistName)
    }

    tracksFound := processSpotifySearchResponse(response)

    return tracksFound, nil
}

// Retrieves a Spotify track that matches the given YouTube video title.
// Sacrifices precision for higher chances of returning a result. 
func (c *SpotifyClient) SearchTracksUsingVideoTitle(accessToken, videoTitle string) ([]utils.UnifiedTrackSearchResult, error) {    
    url := c.buildSearchURL("", "", videoTitle, 1, 0)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request %w", err)
    }

    log.Printf("Sending request to Spotify API: %s", url)

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }

    log.Printf("Received Spotify API response with status code: %d", res.StatusCode)
    if res.StatusCode != http.StatusOK {
        bodyBytes, err := io.ReadAll(res.Body)
        if err == nil {
            log.Printf("Spotify API response body: %s", string(bodyBytes))
        }
    }
    
    defer res.Body.Close()

    var response SpotifySearchResponse
    if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("error decoding response: %w", err)
    }

    if len(response.Tracks.Items) == 0 {
        return nil, fmt.Errorf("no tracks found for YouTube video title %s", videoTitle)
    }

    tracksFound := processSpotifySearchResponse(response)

    return tracksFound, nil
}

// Deletes (unfollows) a playlist in the target user's account
func (c *SpotifyClient) DeletePlaylist(accessToken, playlistID string) error {
    url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", playlistID)

    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return fmt.Errorf("error creating request for DeletePlaylist: %w", err)
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("error executing DeletePlaylist request: %w", err)
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
        body, _ := io.ReadAll(res.Body)
        return fmt.Errorf("spotify API error (status code %d): %s", res.StatusCode, string(body))
    }

    return nil
}