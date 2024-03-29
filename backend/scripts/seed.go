package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/brianvoe/gofakeit/v6"
// 	"github.com/roblieblang/luthien/backend/internal/config"
// 	"github.com/roblieblang/luthien/backend/internal/user"
// 	"github.com/roblieblang/luthien/backend/internal/utils"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func main() {
//     // envConfig := utils.LoadENV()

//     // mongoClient := config.DBConnect(envConfig.MongoURI)

// 	// defer func() {
//         // if err := mongoClient.Disconnect(context.Background()); err != nil {
//     //         log.Fatalf("Failed to disconnect MongoDB client: %v", err)
//     //     }
//     // }()

// 	// dropCtx, dropCancel := context.WithTimeout(context.Background(), 5*time.Second)
//     // defer dropCancel()

//     // usersCollection := mongoClient.Database(envConfig.DatabaseName).Collection("users")

// 	// if err := usersCollection.Drop(dropCtx); err != nil {
// 	// 	log.Fatalf("Failed to drop users collection: %v", err)
// 	// }

// 	// numUsers := 10
// 	// users := make([]interface{}, numUsers)

// 	// for i := 0; i < numUsers; i++ {
// 	// 	gofakeit.Seed(0)
// 	// 	users[i] = user.User {
// 	// 		ID:        primitive.NewObjectID(),
// 	// 		Username:  gofakeit.Username(),
// 	// 		Email:     gofakeit.Email(),
// 	// 		Password:  gofakeit.Password(true, true, true, false, false, 12),
// 	// 		FirstName: gofakeit.FirstName(),
// 	// 		LastName:  gofakeit.LastName(),
// 	// 		CreatedAt: time.Now(),
// 	// 		UpdatedAt: time.Now(),
// 	// 	}
// 	// }

//     // res, err := usersCollection.InsertMany(dropCtx, users)
//     // if err != nil {
//     //     log.Fatalf("Failed to insert users: %v", err)
//     // }

//     // log.Printf("Database seeded successfully: %v", res)
// }