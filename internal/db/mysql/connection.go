package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

// NewConnection initializes a new MySQL database connection
func NewConnection() (*sql.DB, error) {
	// Get MySQL credentials from environment variables
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DB")

	// Validate required environment variables
	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		return nil, fmt.Errorf("❌ missing required MySQL environment variables")
	}

	// Format MySQL connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	// Open MySQL connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to open MySQL connection: %w", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("❌ failed to ping MySQL: %w", err)
	}

	log.Println("✅ Connected to MySQL successfully!")
	return db, nil
}
