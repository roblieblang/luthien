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
		log.Printf("\n\n\nmissing something. SHIT!")
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

// Checks Spotify authentication status for a specific user
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

    _, err := h.youTubeService.CreatePlaylist(playlistData.UserID, playlistData.Payload)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "error creating playlist"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully created new playlist"})
}