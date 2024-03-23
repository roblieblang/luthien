package openai

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OpenAIHandler struct {
	openAIService *OpenAIService
}

func NewOpenAIHandler(openAIService *OpenAIService) *OpenAIHandler {
	return &OpenAIHandler{
		openAIService: openAIService,
	}
}

type ExtractArtistAndSongBody struct {
    VideoTitles []string `json:"videoTitles"`
}

func (h *OpenAIHandler) ExtractArtistAndSongFromVideoTitleHandler(c *gin.Context) {
    var requestBody ExtractArtistAndSongBody
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    resp, err := h.openAIService.ExtractArtistAndSongFromVideoTitle(requestBody.VideoTitles)
    if err != nil {
        log.Printf("Error extracting artist and song: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error extracting artist and song title: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully extracted artist and song from video title", "result": resp})
}
