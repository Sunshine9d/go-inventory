package db

import (
	"database/sql"
	"github.com/Sunshine9d/go-inventory/internal/db/mysql"
	"github.com/Sunshine9d/go-inventory/internal/db/postgres"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"log"
)

type Database interface {
	NewConnection() (*sql.DB, error)
}

// GetProductRepository initializes a product repository
func GetProductRepository() (products.Repository, error) {
	dbConn, err := NewConnection()
	if err != nil {
		return nil, err
	}
	return &mysql.MySQLProductRepository{
		DB:                    dbConn.Gorm,
		SQLDB:                 dbConn.SQLDB,
		GormProductRepository: &repository.GormProductRepository{DB: dbConn.Gorm}, // ✅ Inject GORM logic
	}, nil
}

// GetOrderRepository initializes an order repository
func GetOrderRepository() (orders.Repository, error) {
	dbConn, err := NewConnection()
	if err != nil {
		return nil, err
	}
	return &postgres.PostgresOrderRepository{
		DB:                  dbConn.Gorm,
		SQLDB:               dbConn.SQLDB,
		GormOrderRepository: &repository.GormOrderRepository{DB: dbConn.Gorm}, // ✅ Inject GORM logic
	}, nil
}

// MigrateDB handles automatic migrations for GORM
func MigrateDB() {
	dbConn, err := NewConnection()
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	err = dbConn.Gorm.AutoMigrate(&products.Product{}, &orders.Order{})
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ Database migrated successfully!")
}
