package config

import (
	"os"
	"strconv"
)

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// LoadRedisConfig reads Redis settings from .env
func LoadRedisConfig() RedisConfig {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}
