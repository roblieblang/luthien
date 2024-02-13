package spotify

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
    // TODO: troubleshoot empty access token response(?: if it starts disrupting user experience)
    c.JSON(http.StatusOK, gin.H{"redirectURL": "http://localhost:5173/"})
}

// Checks Spotify authentication status for a specific user
// TODO: could we cache the result of this handler?
func (h *SpotifyHandler) CheckAuthHandler(c *gin.Context) {
    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    userMetadata, err := h.spotifyService.Auth0Service.GetUserMetadata(userID) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user metadata"})
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

    if err := h.spotifyService.HandleLogout(userID); err != nil {
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, userProfile)
}

// Handles the retrieval of the current user's Spotify playlists
// TODO: probably doesn't need to be exposed as an endpoint. 
    // will just need a single endpoint that gets the playlists, 
    // each individual playlist therein, and each individual playlist's items
func(h *SpotifyHandler) GetCurrentUserPlaylistsHandler(c *gin.Context) {
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
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter is required"})
        return
    }
    
    userPlaylists, err := h.spotifyService.GetCurrentUserPlaylists(userID, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, userPlaylists)
}