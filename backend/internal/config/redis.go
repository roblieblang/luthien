package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	var opt *redis.Options
	var err error

	if os.Getenv("GIN_MODE") == "release" {
		opt, err = redis.ParseURL(os.Getenv("REDIS_URL"))
        if err != nil {
            log.Fatalf("Error parsing Redis URL: %v", err)
        }
	} else {
		opt = &redis.Options{
			Addr: addr,
			Password: password,
			DB: db, 
		}
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}
	
	return client
}