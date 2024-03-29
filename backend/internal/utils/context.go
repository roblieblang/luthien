package utils

import (
	"github.com/redis/go-redis/v9"
	// "go.mongodb.org/mongo-driver/mongo"
)

type AppContext struct {
    RedisClient *redis.Client
    // MongoClient *mongo.Client
    EnvConfig   *EnvConfig
}