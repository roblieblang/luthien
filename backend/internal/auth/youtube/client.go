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
        playlists = append(playlists, Playlist{
            ID:          item.Id,
            Title:       item.Snippet.Title,
            Description: item.Snippet.Description,
            ImageURL:    item.Snippet.Thumbnails.Standard.Url,
            VideosCount: item.ContentDetails.ItemCount,
        })
    }

    return YouTubePlaylistsResponse{Playlists: playlists}, nil
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

    call := service.PlaylistItems.List([]string{"snippet", "contentDetails"}).PlaylistId(playlistID)
    resp, err := call.Do()
    if err != nil {
        return YouTubePlaylistItemsResponse{}, fmt.Errorf("error making API call: %v", err)
    }

    var items []PlaylistItem
    for _, item := range resp.Items {
        fmt.Printf("playlist item: %v\n", item)
        var thumbnailURL string
        // Check if the thumbnail and the standard quality thumbnail exist before accessing the URL
        if item.Snippet.Thumbnails != nil && item.Snippet.Thumbnails.Standard != nil {
            thumbnailURL = item.Snippet.Thumbnails.Standard.Url
        }
        items = append(items, PlaylistItem{
            ID:          item.Id,
            Title:       item.Snippet.Title,
            Description: item.Snippet.Description,
            ThumbnailURL: thumbnailURL, // Use the safe variable instead
            VideoID:     item.ContentDetails.VideoId,
        })
    }

    return YouTubePlaylistItemsResponse{Items: items}, nil
}


// TODO: implement add/create/insert playlist

// TODO: implement add/insert a playlist item

// TODO: implement search for a video