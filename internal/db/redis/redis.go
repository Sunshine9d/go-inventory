package redisdb

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// RedisClient is the global Redis client instance
var RedisClient *redis.Client

// Initialize Redis connection
func InitRedis(addr, password string, db int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test Redis connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
}

// Set key-value pair with TTL
func SetCache(key string, value string, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get cached value by key
func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

