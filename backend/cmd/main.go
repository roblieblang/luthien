package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
	"github.com/roblieblang/luthien/backend/internal/auth/spotify"
	"github.com/roblieblang/luthien/backend/internal/db"
	"github.com/roblieblang/luthien/backend/internal/user"
	"github.com/roblieblang/luthien/backend/internal/utils"
)

func main() {
    // Load environment variables
    envConfig := utils.LoadENV()

    // Connect to Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
        Password: "",
        DB: 0,
    })
    _, err := rdb.Ping(context.Background()).Result()
    if err != nil {
        panic(err)
    }

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
        AllowOrigins:     []string{"http://localhost:8080", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
    }))

    userDAO := user.NewDAO(mongoClient, envConfig.DatabaseName, "users")
    userService := user.NewUserService(userDAO)
    userHandler := user.NewUserHandler(userService)

    router.POST("/users", userHandler.CreateUser)
    router.GET("/users", userHandler.GetAllUsers)
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
	if err := router.Run(":" + envConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
