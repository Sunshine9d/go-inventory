package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds all configuration values
type AppConfig struct {
	AppName     string
	Environment string
	ServerPort  string
	DebugMode   bool
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
}

// Global Config variable
var Config AppConfig

// LoadConfig initializes the AppConfig struct
func LoadConfig() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Parse environment variables into AppConfig struct
	Config = AppConfig{
		AppName:     getEnv("APP_NAME", "GoInventory"),
		Environment: getEnv("ENVIRONMENT", "development"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DebugMode:   getEnvAsBool("DEBUG_MODE", false),
		DatabaseURL: getEnv("DATABASE_URL", "mysql://user:password@tcp(localhost:3306)/dbname"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretkey"),
	}
}

// Helper function to get environment variable or fallback to default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as a boolean
func getEnvAsBool(name string, defaultValue bool) bool {
	valueStr := getEnv(name, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
