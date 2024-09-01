package utils

import (
	"github.com/redis/go-redis/v9"
)

type AppContext struct {
	RedisClient *redis.Client
	EnvConfig   *EnvConfig
}
