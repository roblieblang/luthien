package utils

import (
	"github.com/joho/godotenv"

	"log"
	"os"
)

type EnvConfig struct {
    MongoURI            string
    Port                string
    DatabaseName        string
    SpotifyClientID     string
    SpotifyClientSecret string
    SpotifyRedirectURI         string
}

// Load the necessary ENV values
func LoadENV() *EnvConfig {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
    return &EnvConfig{
        MongoURI:            os.Getenv("MONGO_URI"),
        Port:                defaultVal(os.Getenv("PORT"), "8080"),
        DatabaseName:        os.Getenv("MONGO_DB_NAME"),
        SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
        SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
        SpotifyRedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
    }
}

func defaultVal(val, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}