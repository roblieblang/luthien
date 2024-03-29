package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
    RedisAddr                   string
    // MongoURI                    string
    Port                        string
    // DatabaseName                string
    SpotifyClientID             string
    SpotifyClientSecret         string
    SpotifyRedirectURI          string
    GoogleClientID              string
    GoogleClientSecret          string
    GoogleRedirectURI           string
    Auth0ManagementClientID     string
    Auth0ManagementClientSecret string
    Auth0Domain                 string
    OpenAIAPIKey                string
    GinMode                     string
}

// Load the necessary ENV values
func LoadENV() *EnvConfig {
    // if os.Getenv("MONGO_URI") == "" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Print("No .env file found or environment variables not set in Docker")
    //     }
    // }

    if os.Getenv("K_SERVICE") == "" {
        if err := godotenv.Load(); err != nil {
            log.Print("No .env file found, assuming cloud environment variables are set")
        }
    }

    return &EnvConfig{
        RedisAddr:                      os.Getenv("REDIS_ADDR"),
        // MongoURI:                       os.Getenv("MONGO_URI"),
        Port:                           defaultVal(os.Getenv("PORT"), "8080"),
        // DatabaseName:                   os.Getenv("MONGO_DB_NAME"),
        SpotifyClientID:                os.Getenv("SPOTIFY_CLIENT_ID"),
        SpotifyClientSecret:            os.Getenv("SPOTIFY_CLIENT_SECRET"),
        SpotifyRedirectURI:             os.Getenv("SPOTIFY_REDIRECT_URI"),
        GoogleClientID:                 os.Getenv("GOOGLE_CLIENT_ID"),
        GoogleClientSecret:             os.Getenv("GOOGLE_CLIENT_SECRET"),
        GoogleRedirectURI:              os.Getenv("GOOGLE_REDIRECT_URI"),
        Auth0ManagementClientID:        os.Getenv("AUTH0_MANAGEMENT_CLIENT_ID"),
        Auth0ManagementClientSecret:    os.Getenv("AUTH0_MANAGEMENT_CLIENT_SECRET"),
        Auth0Domain:                    os.Getenv("AUTH0_DOMAIN"),
        OpenAIAPIKey:                   os.Getenv("OPENAI_API_KEY"),
        GinMode:                        os.Getenv("GIN_MODE"),
    }
}

func defaultVal(val, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}