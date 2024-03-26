package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/roblieblang/luthien/backend/internal/utils"
	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
    AppContext *utils.AppContext
}

func NewOpenAIClient(appCtx *utils.AppContext) *OpenAIClient {
    return &OpenAIClient{
        AppContext: appCtx,
    }
}

type ArtistSongPair struct {
    ArtistName string `json:"artistName"`
    SongTitle  string `json:"songTitle"`
}

// TODO: handle token limits and API quotas 
// See https://github.com/pkoukk/tiktoken-go#counting-tokens-for-chat-api-calls for token counting

// Prompts the OpenAI API with a list of video titles from which it will extract artist names and song titles
func (c *OpenAIClient) ExtractArtistAndSongFromVideoTitle(videoTitles []string) ([]ArtistSongPair, error) {
	key := c.AppContext.EnvConfig.OpenAIAPIKey
	if key == "" {
		return nil, fmt.Errorf("OpenAI API Key is empty")
	}

	client := openai.NewClient(key)

	backtick := "`"
	prompt := `For each of the following YouTube video titles, use your general knowledge to identify 
	and return an array of objects where each object contains the artist's name under 'artistName' 
	and the song title which is most likely to exist on Spotify under 'songTitle'. 
	If the artist's name cannot be confidently inferred or is not known to you, leave 'artistName' blank.
	Ensure there is no additional text in the response as I will be parsing your response 
	directly into a slice of structs of the following format: 
	type ArtistSongPair struct {
		ArtistName string ` + backtick + `json:"artistName"` + backtick + `
		SongTitle  string ` + backtick + `json:"songTitle"` + backtick + `
	}`

	log.Printf("videoTitles to extract from: %v", videoTitles)
	for _, title := range videoTitles {
		prompt += fmt.Sprintf("\n%s", title)
	}
	log.Printf("Prompt sent to OpenAI: %s", prompt)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Temperature: math.SmallestNonzeroFloat32,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error received during OpenAI chat completion request: %v", err)
	}

	log.Printf("Response from OpenAI: %v", resp)
	responseContent := resp.Choices[0].Message.Content

	log.Printf("Attempting to parse JSON: %s", responseContent)

	var artistSongPairs []ArtistSongPair
	err = json.Unmarshal([]byte(responseContent), &artistSongPairs)
	if err != nil {
		log.Printf("Error parsing OpenAI response into struct: %v", err)
		return nil, fmt.Errorf("error parsing OpenAI response into struct: %v", err)
	}
	return artistSongPairs, nil
} 