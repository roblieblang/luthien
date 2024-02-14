package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien/backend/internal/auth/auth0"
	"github.com/roblieblang/luthien/backend/internal/auth/spotify"
	"github.com/roblieblang/luthien/backend/internal/config"
	"github.com/roblieblang/luthien/backend/internal/user"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

func main() {
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

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:8080", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
    }))


    // Users
    userDAO := user.NewDAO(appCtx.MongoClient, appCtx.EnvConfig.DatabaseName, "users")
    userService := user.NewUserService(userDAO)
    userHandler := user.NewUserHandler(userService)

    router.POST("/users", userHandler.CreateUser)
    router.GET("/users", userHandler.GetAllUsers)
    router.GET("/users/:id", userHandler.GetUser)

    // Auth0
    auth0Client := auth0.NewAuth0Client(appCtx)
    auth0Service := auth0.NewAuth0Service(auth0Client, appCtx)

    // Spotify
    spotifyClient := spotify.NewSpotifyClient(appCtx)
    spotifyService := spotify.NewSpotifyService(spotifyClient, auth0Service, appCtx)
    spotifyHandler := spotify.NewSpotifyHandler(spotifyService)

    router.GET("/auth/spotify/login", spotifyHandler.LoginHandler)
    router.POST("/auth/spotify/callback", spotifyHandler.CallbackHandler)
    router.POST("/auth/spotify/logout", spotifyHandler.LogoutHandler)
    router.GET("/auth/spotify/check-auth", spotifyHandler.CheckAuthHandler)
    router.GET("/spotify/current-profile", spotifyHandler.GetCurrentUserProfileHandler)
    router.GET("/spotify/current-user-playlists", spotifyHandler.GetCurrentUserPlaylistsHandler)
    router.GET("/spotify/playlist-tracks", spotifyHandler.GetPlaylistTracksHandler)
    router.POST("/spotify/create-playlist", spotifyHandler.CreatePlaylistHandler)

    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to the server!",
        })
    })

	if err := router.Run(":" + envConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
    // TODO: make sure serving over HTTPS in prod
}
