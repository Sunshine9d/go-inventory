package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}
