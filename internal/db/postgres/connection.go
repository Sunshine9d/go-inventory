package postgres

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// NewConnection initializes a new PostgreSQL database connection
func NewConnection() (*sql.DB, *gorm.DB, error) {
	// Get PostgreSQL credentials from environment variables
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("POSTGRES_SSLMODE") // Usually "disable" in local dev

	// Validate required environment variables
	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		return nil, nil, fmt.Errorf("❌ missing required PostgreSQL environment variables")
	}

	// Format PostgreSQL connection string
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)

	// Open PostgreSQL connection
	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("❌ failed to open PostgreSQL connection: %w", err)
	}

	// Verify connection
	if err = sqlDB.Ping(); err != nil {
		return nil, nil, fmt.Errorf("❌ failed to ping PostgreSQL: %w", err)
	}
	gormDB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	log.Println("✅ Connected to PostgreSQL successfully!")
	return sqlDB, gormDB, nil
}
