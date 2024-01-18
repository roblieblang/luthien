package main

import (
    "context"
    "log"
    "os"

    "github.com/roblieblang/luthien-core-server/internal/db"
    "github.com/roblieblang/luthien-core-server/internal/dao"
    "github.com/roblieblang/luthien-core-server/internal/services"
    "github.com/roblieblang/luthien-core-server/internal/controllers"


    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    // "go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
    uri := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    client, err := db.Connect(uri, ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to connect and ping database: %v", err)
	}

    userDAO := dao.NewUserDAO(client, "luthien_db", "users")
    userService := services.NewUserService(userDAO)
    userController := controllers.NewUserController(userService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

    router := gin.Default()

    router.POST("/users", userController.CreateUser)
    router.GET("/users/:id", userController.GetUser)

    // router.GET("/", func(c *gin.Context) {
    //     c.JSON(200, gin.H{
    //         "message": "Welcome to the server!",
    //     })
    // })

	if err := router.Run("127.0.0.1:"+port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
