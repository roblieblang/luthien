package youtube

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

type YouTubeHandler struct {
	youTubeService *YouTubeService
}

func NewYouTubeHandler(youTubeService *YouTubeService) *YouTubeHandler {
	return &YouTubeHandler{
		youTubeService: youTubeService,
	}
}

// Sends the session ID and redirect auth URL to the frontend
func (h *YouTubeHandler) LoginHandler(c *gin.Context) {
    authURL, sessionID, err := h.youTubeService.StartLoginFlow()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"authURL": authURL, "sessionID": sessionID})
} 

// Handles a Google de-authentication
func (h *YouTubeHandler) LogoutHandler(c *gin.Context) {
    log.Printf("Inside LogoutHandler")
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

    clearParams := utils.ClearTokensParams{
        Party: "google", 
        UserID: userID, 
        AppCtx: *h.youTubeService.YouTubeClient.AppContext,
    }
    if err := utils.HandleLogout(h.youTubeService.Auth0Service, clearParams); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// After user authorizes the application, Google redirects to a specified callback URL
func (h *YouTubeHandler) CallbackHandler(c *gin.Context) {
    var req struct {
        Code   string `json:"code"`
        UserID string `json:"userID"`
        SessionID string `json:"sessionID"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
	log.Printf("\ncode: %s, userID: %s, sessionID: %s\n", req.Code, req.UserID, req.SessionID)
    if req.Code == "" || req.UserID == "" || req.SessionID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
        return
    }

    err := h.youTubeService.HandleCallback(req.Code, req.UserID, req.SessionID)
    if err != nil {
        log.Printf("Error handling callback: %v\n", err)
        statusCode := http.StatusInternalServerError
        if strings.Contains(err.Error(), "empty access token") {
            statusCode = http.StatusBadRequest
        }
        c.JSON(statusCode, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"redirectURL": "http://localhost:5173/"})
}

// Checks YouTube authentication status for a specific user
// TODO: could we cache the result of this handler?
func (h *YouTubeHandler) CheckAuthHandler(c *gin.Context) {
    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    userMetadata, err := h.youTubeService.Auth0Service.GetUserMetadata(userID) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
    }

    c.JSON(http.StatusOK, gin.H{"isAuthenticated": userMetadata.AppMetadata.AuthenticatedWithGoogle})
}

// Handles the retrieval of the current user's YouTube playlists
func (h *YouTubeHandler) GetCurrentUserPlaylistsHandler(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
		return
	}

	userPlaylists, err := h.youTubeService.GetCurrentUserPlaylists(userID)
	if err != nil {
		log.Printf("Error retrieving YouTube playlists: %v", err)
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with YouTube (Google)."})
            return
        }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user playlists"})
		return
	}

	c.JSON(http.StatusOK, userPlaylists)
}

// Handles the retrieval of the current user's Spotify playlists
func (h *YouTubeHandler) GetPlaylistItemsHandler(c *gin.Context) {
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

	userPlaylists, err := h.youTubeService.GetPlaylistItems(userID, playlistID)
	if err != nil {
		log.Printf("Error retrieving YouTube playlist items: %v", err)
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with YouTube (Google)."})
            return
        }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playlist items"})
		return
	}

	c.JSON(http.StatusOK, userPlaylists)
}

type CreatePlaylistBody struct {
    UserID        string                `json:"userId"`
    Payload       CreatePlaylistPayload `json:"payload"`
}

// Handles the creation of a new playlist 
func(h *YouTubeHandler) CreatePlaylistHandler(c *gin.Context) {
    var playlistData CreatePlaylistBody
    if err := c.BindJSON(&playlistData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    if playlistData.UserID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user Id is required to create a new playlist"})
        return
    }
    createdPlaylist, err := h.youTubeService.CreatePlaylist(playlistData.UserID, playlistData.Payload)
    if err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with YouTube (Google)."})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": "error creating playlist"})
        return
    }

    c.JSON(http.StatusOK, createdPlaylist)
}

type AddItemsToPlaylistBody struct {
    UserID  string                     `json:"userId"`
    Payload AddItemsToPlaylistPayload  `json:"payload"`
}

// Handles the insertion of multiple items into a YouTube playlist
func(h *YouTubeHandler) AddItemsToPlaylistHandler(c *gin.Context) {
    var addItemsData AddItemsToPlaylistBody
    if err := c.BindJSON(&addItemsData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }
    if addItemsData.UserID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user Id is required to create a new playlist"})
        return
    }

    if err := h.youTubeService.AddItemsToPlaylist(addItemsData.UserID, addItemsData.Payload); err != nil {
        errMsg := err.Error()
        if strings.Contains(errMsg, "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with YouTube (Google)."})
            return
        }
        log.Printf("Error adding items to YouTube Playlist: %s", errMsg)
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg, "message": "Error adding items to YouTube Playlist"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully added items to playlist"})
}

// Handles the retrieval of videos that match the given artist name and song title
func (h *YouTubeHandler) SearchVideosHandler(c *gin.Context) {
    userID := c.Query("userID")
    artistName := c.Query("artistName")
    songTitle := c.Query("songTitle")

    if userID == "" {
        log.Printf("userID missing from query parameters")
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }

    // Check if both artistName and songTitle are empty; if so, return an error.
    if artistName == "" && songTitle == "" {
        log.Printf("Both artistName and songTitle missing from query parameters")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Either artist name or song title is required"})
        return
    }

    // Assuming `SearchVideos` method has been adjusted to handle queries with either artistName, songTitle, or both.
    searchResponse, err := h.youTubeService.SearchVideos(userID, artistName, songTitle)
    if err != nil {
        if strings.Contains(err.Error(), "reauthentication required") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication_required", "message": "Please reauthenticate with YouTube (Google)."})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, searchResponse)
}
