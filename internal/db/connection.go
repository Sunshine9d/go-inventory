package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// DBConnection holds both native SQL and GORM DB connections
type DBConnection struct {
	SQLDB *sql.DB
	Gorm  *gorm.DB
}

// NewConnection initializes a database connection based on DB_TYPE
func NewConnection() (*DBConnection, error) {
	dbType := os.Getenv("DB_TYPE")
	// Read environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	var sqlDB *sql.DB
	var gormDB *gorm.DB
	var err error
	switch dbType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbname) // ✅ Add dbname here

		sqlDB, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
		}

		gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize GORM for MySQL: %w", err)
		}

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		fmt.Println(dsn)
		sqlDB, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}

		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize GORM for PostgreSQL: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported database type. Set DB_TYPE to 'mysql' or 'postgres'")
	}

	// Verify connection
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully:", dbType)
	return &DBConnection{SQLDB: sqlDB, Gorm: gormDB}, nil
}
