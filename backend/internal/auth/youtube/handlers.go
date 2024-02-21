package youtube

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
)

type YouTubeHandler struct {
	youTubeService *YouTubeService
}

func NewYouTubeHandler(youTubeService *YouTubeService) *YouTubeHandler {
	return &YouTubeHandler{
		youTubeService: youTubeService,
	}
}

// Handles the retrieval of the current user's Spotify playlists
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