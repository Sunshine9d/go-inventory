package db

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/db/mysql"
	"github.com/Sunshine9d/go-inventory/internal/db/postgres"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"gorm.io/gorm"
	"log"
	"os"
)

type Database interface {
	NewConnection() (*sql.DB, error)
}

// GetProductRepository initializes the correct database repository based on DB_TYPE
func GetProductRepository() (products.Repository, error) {
	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "mysql":
		gormDB, sqlDB, err := mysql.NewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
		}
		return &mysql.MySQLProductRepository{SQLDB: sqlDB, DB: gormDB}, nil

	case "postgres":
		gormDB, sqlDB, err := mysql.NewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
		return &postgres.PostgresProductRepository{SQLDB: sqlDB, DB: gormDB}, nil

	default:
		return nil, fmt.Errorf("unsupported database type. Set DB_TYPE to 'mysql' or 'postgres'")
	}
}

// GetOrderRepository initializes the correct database repository based on DB_TYPE
func GetOrderRepository() (orders.Repository, error) {
	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "mysql":
		dbConn, err := mysql.NewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
		}
		return &mysql.MySQLOrderRepository{DB: dbConn}, nil

	case "postgres":
		dbConn, err := postgres.NewConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
		return &postgres.PostgresOrderRepository{DB: dbConn}, nil

	default:
		return nil, fmt.Errorf("unsupported database type. Set DB_TYPE to 'mysql' or 'postgres'")
	}
}

// NewGORMConnection initializes GORM based on environment variables
func NewGORMConnection() (*gorm.DB, error) {
	dbType := os.Getenv("DB_TYPE")
	var dsn string

	switch dbType {
	case "mysql":
		dsn = os.Getenv("MYSQL_DSN") // Example: "user:password@tcp(localhost:3306)/inventory?charset=utf8mb4&parseTime=True&loc=Local"
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn = os.Getenv("POSTGRES_DSN") // Example: "host=localhost user=postgres password=pass dbname=inventory port=5432 sslmode=disable"
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// MigrateDB handles automatic migrations for GORM
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&products.Product{}, &orders.Order{}, &orders.OrderItem{})
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	log.Println("✅ Database migrated successfully!")
}
