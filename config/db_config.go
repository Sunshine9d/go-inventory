package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error

	// Choose the database driver based on ENV variable
	dbType := os.Getenv("DB_TYPE") // "mysql" or "postgres"
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	switch dbType {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatal("Unsupported DB_TYPE. Use 'mysql' or 'postgres'.")
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connected successfully!")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
