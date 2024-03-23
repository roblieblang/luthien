package openai

type OpenAIService struct {
	OpenAIClient *OpenAIClient
}

func NewOpenAIService(openAIClient *OpenAIClient) *OpenAIService {
	return &OpenAIService{
		OpenAIClient: openAIClient,
	}
}

func (s *OpenAIService) ExtractArtistAndSongFromVideoTitle(videoTitles []string) ([]ArtistSongPair, error) {
	return s.OpenAIClient.ExtractArtistAndSongFromVideoTitle(videoTitles)
}