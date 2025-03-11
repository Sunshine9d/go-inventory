package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// NewConnection initializes a new MySQL database connection
func NewConnection() (*sql.DB, *gorm.DB, error) {
	// Get MySQL credentials from environment variables
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DB")
	fmt.Println(user, password, host, port, dbname)
	// Validate required environment variables
	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		return nil, nil, fmt.Errorf("❌ missing required MySQL environment variables")
	}

	// Format MySQL connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	// Open MySQL connection
	sqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("❌ failed to open MySQL connection: %w", err)
	}

	// Verify connection
	if err = sqlDB.Ping(); err != nil {
		return nil, nil, fmt.Errorf("❌ failed to ping MySQL: %w", err)
	}

	gormDB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	log.Println("✅ Connected to MySQL successfully!")
	return sqlDB, gormDB, nil
}
