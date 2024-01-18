package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/roblieblang/luthien-core-server/internal/auth/spotify"
    "github.com/roblieblang/luthien-core-server/internal/user"
	"github.com/roblieblang/luthien-core-server/internal/db"
    "github.com/roblieblang/luthien-core-server/internal/utils"

)

func main() {
    // Load environment variables
    envConfig := utils.LoadENV()


    // Connect to MongoDB
    mongoClient := db.Connect(envConfig.MongoURI)
    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            log.Fatalf("Failed to disconnect MongoDB client: %v", err)
        }
    }()


    // Create the router and assign routes
    router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
    }))

    userDAO := user.NewDAO(mongoClient, envConfig.DatabaseName, "users")
    userService := user.NewUserService(userDAO)
    userHandler := user.NewUserHandler(userService)

    router.POST("/users", userHandler.CreateUser)
    router.GET("/users/:id", userHandler.GetUser)


    router.GET("/auth/spotify/login", func(c *gin.Context) {
        spotify.LoginHandler(c, envConfig.SpotifyClientID, envConfig.SpotifyRedirectURI)
    })
    router.GET("/auth/spotify/callback", func(c *gin.Context) {
        spotify.CallbackHandler(c, envConfig.SpotifyClientID, envConfig.SpotifyRedirectURI)
    })
    // TODO: not sure if that's the right route...
    router.GET("/auth/spotify", spotify.RefreshTokenHandler)

    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to the server!",
        })
    })


    // Run the server
	if err := router.Run("127.0.0.1:" + envConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
