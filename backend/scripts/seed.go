package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"github.com/roblieblang/luthien-core-server/internal/db"
	"github.com/roblieblang/luthien-core-server/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

    uri := os.Getenv("MONGO_URI")
    databaseName := os.Getenv("MONGO_DB_NAME")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := db.Connect(uri, ctx)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer client.Disconnect(ctx)

    usersCollection := client.Database(databaseName).Collection("users")

	if err := usersCollection.Drop(ctx); err != nil {
		log.Fatalf("Failed to drop users collection: %v", err)
	}

	numUsers := 10
	users := make([]interface{}, numUsers)

	for i := 0; i < numUsers; i++ {
		gofakeit.Seed(0)
		users[i] = models.User {
			ID:        primitive.NewObjectID(),
			Username:  gofakeit.Username(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, false, false, 12),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

    _, err = usersCollection.InsertMany(ctx, users)
    if err != nil {
        log.Fatalf("Failed to insert users: %v", err)
    }

    log.Println("Database seeded successfully")
}
