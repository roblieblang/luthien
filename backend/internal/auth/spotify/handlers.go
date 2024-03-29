package spotify

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

// Parses incoming HTTP requests for parameters, payloads, headers
// Performs inital validation of the request before passing it onto service layer
// Formats responses and errors from service layer into HTTP format

type SpotifyHandler struct {
    spotifyService *SpotifyService
}

func NewSpotifyHandler(spotifyService *SpotifyService) *SpotifyHandler {
    return &SpotifyHandler{
        spotifyService: spotifyService,
    }
}

// Sends the session ID and redirect auth URL to the frontend
func (h *SpotifyHandler) LoginHandler(c *gin.Context) {
    authURL, sessionID, err := h.spotifyService.StartLoginFlow()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"authURL": authURL, "sessionID": sessionID})
} 

// Once user authorizes the application, Spotify redirects to a callback URL specified in application settings.
// This handler is called by Spotify, not our own application 
func (h *SpotifyHandler) CallbackHandler(c *gin.Context) {
    var req struct {
        Code   string `json:"code"`
        UserID string `json:"userID"`
        SessionID string `json:"sessionID"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    if req.Code == "" || req.UserID == "" || req.SessionID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
        return
    }

    err := h.spotifyService.HandleCallback(req.Code, req.UserID, req.SessionID)
    if err != nil {
        log.Printf("Error handling callback: %v\n", err)
        statusCode := http.StatusInternalServerError
        if strings.Contains(err.Error(), "empty access token") {
            statusCode = http.StatusBadRequest
        }
        c.JSON(statusCode, gin.H{"error": err.Error()})
        return
    }

    redirectURL := "http://localhost:5173/"
    if os.Getenv("GIN_MODE") == "release" {
        redirectURL = os.Getenv("DEPLOYED_UI_URL")
    }
    
    c.JSON(http.StatusOK, gin.H{"redirectURL": redirectURL})
}

// Checks Spotify authentication status for a specific user
func (h *SpotifyHandler) CheckAuthHandler(c *gin.Context) {
    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    userMetadata, err := h.spotifyService.Auth0Service.GetUserMetadata(userID) 
    if err != nil {
        log.Printf("Error getting Auth0 user metadata: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
    }

    c.JSON(http.StatusOK, gin.H{"isAuthenticated": userMetadata.AppMetadata.AuthenticatedWithSpotify})
}

// Handles a Spotify logout(de-authentication)
func (h *SpotifyHandler) LogoutHandler(c *gin.Context) {
    var req struct {
        UserID string `json:"userID"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }  
    
    userID := req.UserID
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    clearTokenParams := utils.ClearTokensParams{
        Party: "spotify", 
        UserID: userID,
        AppCtx: *h.spotifyService.AppContext,
    }
    if err := utils.HandleLogout(h.spotifyService.Auth0Service, clearTokenParams); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Handles the retrieval of the current user's Spotify profile data
func (h *SpotifyHandler) GetCurrentUserProfileHandler(c *gin.Context) {
    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }

    userProfile, err := h.spotifyService.GetCurrentUserProfile(userID) 
    if err != nil {
        log.Printf("error retrieving Spotify profile: %v", err)
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, userProfile)
}

// Handles the retrieval of the current user's Spotify playlists
func(h *SpotifyHandler) GetCurrentUserPlaylistsHandler(c *gin.Context) {
    defaultOffset := 0

    offsetStr := c.DefaultQuery("offset", strconv.Itoa(defaultOffset))
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        offset = defaultOffset
    }

    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }
    
    userPlaylists, err := h.spotifyService.GetCurrentUserPlaylists(userID, offset)
    if err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, userPlaylists)
}

// Handles the retrieval of a single playlist's tracks
func(h *SpotifyHandler) GetPlaylistTracksHandler(c *gin.Context) {
    // defaultLimit := 20
    // defaultOffset := 0

    // limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
    // limit, err := strconv.Atoi(limitStr)
    // if err != nil {
    //     limit = defaultLimit
    // }

    // offsetStr := c.DefaultQuery("offset", strconv.Itoa(defaultOffset))
    // offset, err := strconv.Atoi(offsetStr)
    // if err != nil {
    //     offset = defaultOffset
    // }

    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }

    playlistID := c.Query("playlistID")
    if playlistID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "playlistID query parameter is required"})
        return
    }
    
    playlistTracks, err := h.spotifyService.GetPlaylistTracks(userID, playlistID)
    if err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, playlistTracks)
}

// Handles the creation of a new playlist 
func(h *SpotifyHandler) CreatePlaylistHandler(c *gin.Context) {
    var playlistData CreatePlaylistBody
    if err := c.BindJSON(&playlistData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    newPlaylistID, err := h.spotifyService.CreatePlaylist(playlistData.UserID, playlistData.SpotifyUserID, playlistData.Payload)
    if err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": "error creating playlist"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully created new playlist", "newPlaylistId": newPlaylistID})
}

type AddItemsToPlaylistBody struct {
    UserID     string                    `json:"userId"`
    PlaylistID string                    `json:"spotifyPlaylistId"`
    Payload    AddItemsToPlaylistPayload `json:"payload"`
}

// Handles the insertion of items into an existing playlist 
func(h *SpotifyHandler) AddItemsToPlaylistHandler(c *gin.Context) {
    var playlistItemsData AddItemsToPlaylistBody
    if err := c.BindJSON(&playlistItemsData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    err := h.spotifyService.AddItemsToPlaylist(playlistItemsData.UserID, playlistItemsData.PlaylistID, playlistItemsData.Payload)
    if err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error adding items to playlist: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully add items to playlist with ID: %s", playlistItemsData.PlaylistID)})
}

// Handles the retrieval of track URI given an artist name and track title
func(h *SpotifyHandler) SearchTracksUsingArtistAndTrackhandler(c *gin.Context) {
    defaultLimit := 20
    defaultOffset := 0

    limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        limit = defaultLimit
    }

    offsetStr := c.DefaultQuery("offset", strconv.Itoa(defaultOffset))
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        offset = defaultOffset
    }

    userID := c.Query("userID")
    if userID == "" {
        log.Printf("userID missing from query parameters")
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }

    trackTitle := c.Query("trackTitle")
    artistName := c.Query("artistName")

    if artistName == "" && trackTitle == "" {
        log.Printf("Either trackTitle or artistName must be provided")
        c.JSON(http.StatusBadRequest, gin.H{"error": "either trackTitle or artistName query parameter is required"})
        return
    }
    
    tracksFound, err := h.spotifyService.SearchTracksUsingArtistAndTrack(userID, artistName, trackTitle, limit, offset)
    if err != nil {
        log.Printf("Search error: %v", err)
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        } else if strings.Contains(err.Error(), "no tracks found") {
            c.JSON(http.StatusNotFound, gin.H{"error": "No tracks found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        }
        return
    }

    c.JSON(http.StatusOK, tracksFound)
}

// Handles the retrieval of track URI given a video title
func(h *SpotifyHandler) SearchTracksUsingVideoTitleHandler(c *gin.Context) {
    userID := c.Query("userID")
    if userID == "" {
        log.Printf("userID missing from query parameters")
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }

    videoTitle := c.Query("videoTitle")

    if videoTitle == "" {
        log.Printf("videoTitle must be provided")
        c.JSON(http.StatusBadRequest, gin.H{"error": "videoTitle query parameter is required"})
        return
    }
    
    tracksFound, err := h.spotifyService.SearchTracksUsingVideoTitle(userID, videoTitle)
    if err != nil {
        log.Printf("Search error: %v", err)
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        } else if strings.Contains(err.Error(), "no tracks found") {
            c.JSON(http.StatusNotFound, gin.H{"error": "No tracks found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        }
        return
    }

    c.JSON(http.StatusOK, tracksFound)
}

// Handles the deletion of a Spotify playlist
func(h *SpotifyHandler) DeletePlaylistHandler(c *gin.Context) {
    userID := c.Query("userID")
    if userID == "" {
        log.Printf("userID missing from query parameters")
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }

    playlistID := c.Query("playlistID")
    if playlistID == "" {
        log.Printf("playlistID missing from query parameters")
        c.JSON(http.StatusBadRequest, gin.H{"error": "playlistID query parameter is required"})
        return
    }

    if err := h.spotifyService.DeletePlaylist(userID, playlistID); err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with Spotify."})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error deleting playlist: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "playlist deleted successfully"})
}