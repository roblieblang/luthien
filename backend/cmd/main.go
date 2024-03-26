package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/auth/openai"
	"github.com/roblieblang/luthien/backend/internal/auth/spotify"
	"github.com/roblieblang/luthien/backend/internal/auth/youtube"
	"github.com/roblieblang/luthien/backend/internal/config"
	"github.com/roblieblang/luthien/backend/internal/user"
	"github.com/roblieblang/luthien/backend/internal/utils"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		log.Printf("[%s] %s %s %d in %v", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}

func main() {
    log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/server.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true, // compress rolled back files
	})

    envConfig := utils.LoadENV()

    redisClient := config.NewRedisClient(envConfig.RedisAddr, "", 0)

    mongoClient:= config.DBConnect(envConfig.MongoURI)
    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            log.Fatalf("Failed to disconnect MongoDB client: %v", err)
        }
    }()

    appCtx := &utils.AppContext{
        EnvConfig:   envConfig,
        RedisClient: redisClient,
        MongoClient: mongoClient,
    }

    router := gin.Default()

    router.Use(LoggerMiddleware())

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:8080", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
    }))


    // User setup
    userDAO := user.NewDAO(appCtx.MongoClient, appCtx.EnvConfig.DatabaseName, "users")
    userService := user.NewUserService(userDAO)
    userHandler := user.NewUserHandler(userService)

    // User data endpoints
    router.POST("/users", userHandler.CreateUser)
    router.GET("/users", userHandler.GetAllUsers)
    router.GET("/users/:id", userHandler.GetUser)

    // Auth0 setup
    auth0Client := auth0.NewAuth0Client(appCtx)
    auth0Service := auth0.NewAuth0Service(auth0Client, appCtx)

    // Spotify setup
    spotifyClient := spotify.NewSpotifyClient(appCtx)
    spotifyService := spotify.NewSpotifyService(spotifyClient, auth0Service, appCtx)
    spotifyHandler := spotify.NewSpotifyHandler(spotifyService)

    // Spotify authentication endpoints
    router.GET("/auth/spotify/login", spotifyHandler.LoginHandler)
    router.POST("/auth/spotify/callback", spotifyHandler.CallbackHandler)
    router.POST("/auth/spotify/logout", spotifyHandler.LogoutHandler)
    router.GET("/auth/spotify/check-auth", spotifyHandler.CheckAuthHandler)

    // Spotify user data endpoints
    router.GET("/spotify/current-profile", spotifyHandler.GetCurrentUserProfileHandler)
    router.GET("/spotify/current-user-playlists", spotifyHandler.GetCurrentUserPlaylistsHandler)
    router.GET("/spotify/playlist-tracks", spotifyHandler.GetPlaylistTracksHandler)
    router.POST("/spotify/create-playlist", spotifyHandler.CreatePlaylistHandler)
    router.POST("/spotify/add-items-to-playlist", spotifyHandler.AddItemsToPlaylistHandler)
    router.GET("/spotify/search-for-track", spotifyHandler.SearchTracksUsingArtistAndTrackhandler)
    router.GET("/spotify/search-using-video", spotifyHandler.SearchTracksUsingVideoTitleHandler)
    router.DELETE("/spotify/delete-playlist", spotifyHandler.DeletePlaylistHandler)

    // YouTube setup
    youTubeClient := youtube.NewYouTubeClient(appCtx)
    youTubeService := youtube.NewYouTubeService(youTubeClient, auth0Service)
    youTubeHandler := youtube.NewYouTubeHandler(youTubeService)

    // Google authentication endpoints
    router.GET("/auth/google/login", youTubeHandler.LoginHandler)
    router.POST("/auth/google/callback", youTubeHandler.CallbackHandler)
    router.POST("/auth/google/logout", youTubeHandler.LogoutHandler)
    router.GET("/auth/google/check-auth", youTubeHandler.CheckAuthHandler)

    // YouTube data endpoints
    router.GET("/youtube/current-user-playlists", youTubeHandler.GetCurrentUserPlaylistsHandler)
    router.GET("/youtube/playlist-tracks", youTubeHandler.GetPlaylistItemsHandler)
    router.POST("/youtube/create-playlist", youTubeHandler.CreatePlaylistHandler)
    router.POST("/youtube/add-items-to-playlist", youTubeHandler.AddItemsToPlaylistHandler)
    router.GET("/youtube/search-for-video", youTubeHandler.SearchVideosHandler)
    router.DELETE("/youtube/delete-playlist", youTubeHandler.DeletePlaylistHandler)

    // OpenAI setup
    openAIClient := openai.NewOpenAIClient(appCtx)
    openAIService := openai.NewOpenAIService(openAIClient)
    openAIHandler := openai.NewOpenAIHandler(openAIService)

    // OpenAI endpoints
    router.POST("/auth/openai/extract-artist-song", openAIHandler.ExtractArtistAndSongFromVideoTitleHandler)


    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to the server!",
        })
    })

	if err := router.Run(":" + envConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
