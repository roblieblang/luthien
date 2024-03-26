package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/roblieblang/luthien/backend/internal/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeClient struct {
    AppContext *utils.AppContext
}

type YouTubePlaylistsResponse struct {
    TotalCount      int        `json:"totalCount"`
    // PrevPageToken   string     `json:"prevPageToken"`
    // NextPageToken   string     `json:"nextPageToken"`
    Playlists       []Playlist `json:"playlists"`
}

type Playlist struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    ImageURL    string `json:"imageUrl"`
    VideosCount int64  `json:"videosCount"`
}

type PlaylistItem struct {
    ID                        string `json:"id"`
    Title                     string `json:"title"`
    Description               string `json:"description"`
    ThumbnailURL              string `json:"thumbnailUrl"`
    VideoID                   string `json:"videoId"`
    VideoOwnerChannelTitle    string `json:"videoOwnerChannelTitle"`
}

type YouTubePlaylistItemsResponse struct {
    Items []PlaylistItem `json:"items"`
}

func NewYouTubeClient(appCtx *utils.AppContext) *YouTubeClient {
    return &YouTubeClient{
        AppContext: appCtx,
    }
}

// Requests a new access token from Google
func (c *YouTubeClient) RequestToken(payload url.Values) (utils.TokenResponse, error) {
    resp, err := http.PostForm("https://oauth2.googleapis.com/token", payload)
    if err != nil {
        log.Printf("error making request for new Google access token: %v", err)
        return utils.TokenResponse{}, err
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("error reading google token response body: %v", err)
        return utils.TokenResponse{}, err
    }
    bodyContent := string(bodyBytes)
    
    if resp.StatusCode >= 400 {
        log.Printf("google access token request failed: %s, response body: %s", resp.Status, bodyContent)
    }

    var tokenResponse utils.TokenResponse
    if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
        log.Printf("error decoding google token response: %v", err)
        return utils.TokenResponse{}, err
    }

    return tokenResponse, nil
}

// Gets the current user's playlists
func (c *YouTubeClient) GetCurrentUserPlaylists(accessToken string) (YouTubePlaylistsResponse, error) {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)
    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        return YouTubePlaylistsResponse{}, fmt.Errorf("error creating YouTube service: %v", err)
    }

    var playlists []Playlist
    var nextPageToken string
    var totalCount int
    for {
        call := service.Playlists.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(50).PageToken(nextPageToken)
        resp, err := call.Do()
        if err != nil {
            return YouTubePlaylistsResponse{}, fmt.Errorf("error making API call: %v", err)
        }

        for _, item := range resp.Items {
            imageURL := getBestAvailableThumbnailURL(item.Snippet.Thumbnails)

            playlists = append(playlists, Playlist{
                ID:          item.Id,
                Title:       item.Snippet.Title,
                Description: item.Snippet.Description,
                ImageURL:    imageURL,
                VideosCount: item.ContentDetails.ItemCount,
            })
        }

        nextPageToken = resp.NextPageToken
        if nextPageToken == "" {
            totalCount = int(resp.PageInfo.TotalResults)
            break
        }
    }

    return YouTubePlaylistsResponse{
        TotalCount: totalCount,
        Playlists:  playlists,
    }, nil
}




// Gets a playlist's items
func (c *YouTubeClient) GetPlaylistItems(playlistID, accessToken string) (YouTubePlaylistItemsResponse, error) {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)

    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        return YouTubePlaylistItemsResponse{}, fmt.Errorf("error creating YouTube service: %v", err)
    }

    var items []PlaylistItem
    var nextPageToken string

    for {
        call := service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
            PlaylistId(playlistID).MaxResults(50).PageToken(nextPageToken)

        resp, err := call.Do()
        if err != nil {
            return YouTubePlaylistItemsResponse{}, fmt.Errorf("error making API call: %v", err)
        }

        for _, item := range resp.Items {
            thumbnailURL := getBestAvailableThumbnailURL(item.Snippet.Thumbnails)

            items = append(items, PlaylistItem{
                ID:                     item.Id,
                Title:                  item.Snippet.Title,
                Description:            item.Snippet.Description,
                ThumbnailURL:           thumbnailURL,
                VideoID:                item.ContentDetails.VideoId,
                VideoOwnerChannelTitle: item.Snippet.VideoOwnerChannelTitle,
            })
        }

        nextPageToken = resp.NextPageToken
        if nextPageToken == "" {
            break
        }
    }

    return YouTubePlaylistItemsResponse{Items: items}, nil
}


// Helper function to determine the best available thumbnail URL
func getBestAvailableThumbnailURL(thumbnails *youtube.ThumbnailDetails) string {
    if thumbnails == nil {
        return ""
    }
    if thumbnails.Maxres != nil {
        return thumbnails.Maxres.Url
    } else if thumbnails.Standard != nil {
        return thumbnails.Standard.Url
    } else if thumbnails.High != nil {
        return thumbnails.High.Url
    } else if thumbnails.Medium != nil {
        return thumbnails.Medium.Url
    } else if thumbnails.Default != nil {
        return thumbnails.Default.Url
    }
    return ""
}

type CreatePlaylistPayload struct {
    Title           string `json:"title"`
    Description     string `json:"description,omitempty"`
    PrivacyStatus   string `json:"privacyStatus,omitempty"` // "public", "private", or "unlisted"
}

// Creates a new YouTube playlist
func (c *YouTubeClient) CreatePlaylist(accessToken string, payload CreatePlaylistPayload) (*youtube.Playlist, error) {
    if payload.Title == ""{
        log.Printf("Missing playlist title in payload for CreatePlaylist YouTube client")
        return nil, fmt.Errorf("missing playlist title")
    }
    
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)

    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        log.Printf("Error creating new YouTube service: %v", err)
        return nil, fmt.Errorf("error creating YouTube service: %v", err)
    }

    // Create the playlist object to be inserted
    playlist := &youtube.Playlist{
        Snippet: &youtube.PlaylistSnippet{
            Title:       payload.Title,
            Description: payload.Description,
        },
        Status: &youtube.PlaylistStatus{
            PrivacyStatus: payload.PrivacyStatus,
        },
    }

    // Call the YouTube Data API to insert the playlist
    call := service.Playlists.Insert([]string{"snippet", "status"}, playlist)
    createdPlaylist, err := call.Do()
    if err != nil {
        log.Printf("Error creating YouTube playlist: %v", err)
        return nil, fmt.Errorf("error creating YouTube playlist: %v", err)
    }

    return createdPlaylist, nil
}

type AddItemsToPlaylistPayload struct {
    PlaylistID string   `json:"playlistId"`
    VideoIDs   []string `json:"videoIds"`
}

// Adds items to an existing YouTube playlist
func (c *YouTubeClient) AddItemsToPlaylist(accessToken string, payload AddItemsToPlaylistPayload) error {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)
    
    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        return fmt.Errorf("error creating YouTube service: %v", err)
    }

    for _, videoID := range payload.VideoIDs {
        playlistItem := &youtube.PlaylistItem{
            Snippet: &youtube.PlaylistItemSnippet{
                PlaylistId: payload.PlaylistID,
                ResourceId: &youtube.ResourceId{
                    Kind:    "youtube#video",
                    VideoId: videoID,
                },
            },
        }
        call := service.PlaylistItems.Insert([]string{"snippet"}, playlistItem)
        _, err := call.Do()
        if err != nil {
            return fmt.Errorf("error adding item to YouTube playlist: %v", err)
        }
    }

    return nil
}

type YouTubeVideoSearchResponse struct {
    Items []VideoSearchResult `json:"items"`
}

type VideoSearchResult struct {
    ID           VideoID `json:"id"`
    Title        string  `json:"title"`
    Description  string  `json:"description"`
    ThumbnailURL string  `json:"thumbnailUrl"`
    ChannelTitle string  `json:"channelTitle"`
}

type VideoID struct {
    VideoID string `json:"videoId"`
}
// TODO: cache search results (search is expensive)
// Searches for videos on YouTube based on a query
func (c *YouTubeClient) SearchVideos(accessToken, query string,  maxResults int64) ([]utils.UnifiedTrackSearchResult, error) {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)

    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        log.Printf("error creating new YouTube service: %v", err)
        return nil, fmt.Errorf("error creating YouTube service: %v", err)
    }

    call := service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(maxResults).Type("video")
    resp, err := call.Do()
    if err != nil {
        log.Printf("error searching for YouTube video: %v", err)
        return nil, fmt.Errorf("error making API call: %v", err)
    }

    var results []utils.UnifiedTrackSearchResult
    for _, item := range resp.Items {
        thumbnailURL := getBestAvailableThumbnailURL(item.Snippet.Thumbnails)
        result := utils.UnifiedTrackSearchResult{
            ID:           item.Id.VideoId,
            Title:        item.Snippet.Title,
            Artist:       "", 
            Album:        "",
            Thumbnail:    thumbnailURL,
        }
        results = append(results, result)
    }

    return results, nil
}

// Deletes the specified YouTube playlist
func(c *YouTubeClient) DeletePlaylist(accessToken, playlistID string) error {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)

    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        log.Printf("error creating new YouTube service: %v", err)
        return fmt.Errorf("error creating YouTube service: %v", err)
    }

    call := service.Playlists.Delete(playlistID)
    err = call.Do()
    if err != nil {
        return fmt.Errorf("deleting YouTube playlist: %w", err)
    }
    return nil
}