package youtube

import (
	"context"
    "fmt"

    "google.golang.org/api/option"
    "google.golang.org/api/youtube/v3"
	"golang.org/x/oauth2"
	"github.com/roblieblang/luthien/backend/internal/utils"
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
    VideosCount int64  `json:"videosCount"`
}

func NewYouTubeClient(appCtx *utils.AppContext) *YouTubeClient {
    return &YouTubeClient{
        AppContext: appCtx,
    }
}

// TODO: handle Google refresh tokens with oauth2 package (or another method)
// See: https://community.auth0.com/t/store-and-retrieve-google-refresh-token/28973/5
// Gets the current user's playlists
func (c *YouTubeClient) GetCurrentUserPlaylists (accessToken string) (YouTubePlaylistsResponse, error){
	token := &oauth2.Token{AccessToken: accessToken}
	tokenSource := oauth2.StaticTokenSource(token)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	// service, err := youtube.NewService(context.Background(), option.WithTokenSource(nil))
	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
    if err != nil {
        return YouTubePlaylistsResponse{}, fmt.Errorf("error creating YouTube service (Google API Go Client): %v", err)
    }

	call := service.Playlists.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(50)
	call.Header().Set("Authorization", "Bearer "+accessToken)

	resp, err := call.Do()
	if err != nil {
		return YouTubePlaylistsResponse{}, fmt.Errorf("error making Google API call: %v", err)
	}

	var playlists []Playlist
    for _, item := range resp.Items {
        playlists = append(playlists, Playlist{
            ID:          item.Id,
            Title:       item.Snippet.Title,
            VideosCount: item.ContentDetails.ItemCount,
        })
    }

	return YouTubePlaylistsResponse{Playlists: playlists}, nil
}

// TODO: implement get playlist items

// TODO: implement add/create/insert playlist

// TODO: implement add/insert a playlist item

// TODO: implement search for a video