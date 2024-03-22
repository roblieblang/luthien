package openai

import (
	"context"
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

// TODO: handle token limits and API quotas 
// See https://github.com/pkoukk/tiktoken-go#counting-tokens-for-chat-api-calls for token counting

// Prompts the OpenAI API with a list of video titles from which it will extract artist names and song titles
func (c *OpenAIClient) ExtractArtistAndSongFromVideoTitle(videoTitles []string) (string, error) {
	key := c.AppContext.EnvConfig.OpenAIAPIKey
	if key == "" {
		return "", fmt.Errorf("OpenAI API Key is empty")
	}

	client := openai.NewClient(key)
	// TODO: instead of a 2D array, should get an array of objects where songtitle and artistname are clearly identified
	prompt := "Extract and return a 2D array containing the artist name and song title from each of the following video titles, with no additional text in the response:"
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
		return "", fmt.Errorf("error received during OpenAI chat completion request: %v", err)
	}
	
	log.Printf("Response from OpenAI: %v", resp)
	return resp.Choices[0].Message.Content, nil
} 