package youtube

import (
	"context"
	"encoding/json"
	"fmt"
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
	Playlists []Playlist `json:"playlists"`
}

type Playlist struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    ImageURL    string `json:"imageUrl"`
    VideosCount int64  `json:"videosCount"`
}

type PlaylistItem struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    ThumbnailURL string `json:"thumbnailUrl"`
    VideoID     string `json:"videoId"`
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
        return utils.TokenResponse{}, err
    }
    defer resp.Body.Close()

    var tokenResponse utils.TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return utils.TokenResponse{}, err
    }
    return tokenResponse, nil
}

// Gets the current user's playlists
// TODO: pagination
func (c *YouTubeClient) GetCurrentUserPlaylists(accessToken string) (YouTubePlaylistsResponse, error) {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)
    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        return YouTubePlaylistsResponse{}, fmt.Errorf("error creating YouTube service: %v", err)
    }

    call := service.Playlists.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(50)
    resp, err := call.Do()
    if err != nil {
        return YouTubePlaylistsResponse{}, fmt.Errorf("error making API call: %v", err)
    }

    var playlists []Playlist
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

    return YouTubePlaylistsResponse{Playlists: playlists}, nil
}


// Gets a playlist's items
// TODO: pagination
func (c *YouTubeClient) GetPlaylistItems(playlistID, accessToken string) (YouTubePlaylistItemsResponse, error) {
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)

    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        return YouTubePlaylistItemsResponse{}, fmt.Errorf("error creating YouTube service: %v", err)
    }

    call := service.PlaylistItems.List([]string{"snippet", "contentDetails"}).PlaylistId(playlistID)
    resp, err := call.Do()
    if err != nil {
        return YouTubePlaylistItemsResponse{}, fmt.Errorf("error making API call: %v", err)
    }

    var items []PlaylistItem
    for _, item := range resp.Items {
        thumbnailURL := getBestAvailableThumbnailURL(item.Snippet.Thumbnails)

        items = append(items, PlaylistItem{
            ID:             item.Id,
            Title:          item.Snippet.Title,
            Description:    item.Snippet.Description,
            ThumbnailURL:   thumbnailURL, // Use the safe variable instead
            VideoID:        item.ContentDetails.VideoId,
        })
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
    token := &oauth2.Token{AccessToken: accessToken}
    tokenSource := oauth2.StaticTokenSource(token)
    httpClient := oauth2.NewClient(context.Background(), tokenSource)

    service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
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
        return nil, fmt.Errorf("error creating YouTube playlist: %v", err)
    }

    return createdPlaylist, nil
}

type AddItemsToPlaylistPayload struct {
    PlaylistID string   `json:"playlistId"`
    VideoIDs   []string `json:"videoIds"`
}

// Adds items to an existing YouTube playlist
func (c *YouTubeClient) AddItemsToPlaylist (accessToken string, payload AddItemsToPlaylistPayload) error {
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

// TODO: implement search for a video