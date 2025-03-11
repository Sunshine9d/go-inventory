package db

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/internal/db/mysql"
	"github.com/Sunshine9d/go-inventory/internal/db/postgres"
	"github.com/Sunshine9d/go-inventory/internal/orders"
	"github.com/Sunshine9d/go-inventory/internal/products"
	"github.com/Sunshine9d/go-inventory/internal/repository"
	"log"
	"os"
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
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "postgres":
		return &postgres.PostgresProductRepository{
			DB:                    dbConn.Gorm,
			SQLDB:                 dbConn.SQLDB,
			GormProductRepository: &repository.GormProductRepository{DB: dbConn.Gorm}, // ✅ Inject GORM logic
		}, nil
	case "mysql":
		return &mysql.MySQLProductRepository{
			DB:                    dbConn.Gorm,
			SQLDB:                 dbConn.SQLDB,
			GormProductRepository: &repository.GormProductRepository{DB: dbConn.Gorm}, // ✅ Inject GORM logic
		}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// GetOrderRepository initializes an order repository
func GetOrderRepository() (orders.Repository, error) {
	dbConn, err := NewConnection()
	if err != nil {
		return nil, err
	}
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "mysql":
		return &mysql.MySQLOrderRepository{
			DB:                  dbConn.Gorm,
			SQLDB:               dbConn.SQLDB,
			GormOrderRepository: &repository.GormOrderRepository{DB: dbConn.Gorm},
		}, nil
	case "postgres":
		return &postgres.PostgresOrderRepository{
			DB:                  dbConn.Gorm,
			SQLDB:               dbConn.SQLDB,
			GormOrderRepository: &repository.GormOrderRepository{DB: dbConn.Gorm},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

}

// MigrateDB handles automatic migrations for GORM
func MigrateDB() {
	dbConn, err := NewConnection()
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "mysql":
		err = mysql.MigrateMySQL(dbConn.SQLDB)
	case "postgres":
		err = postgres.MigratePostgres(dbConn.SQLDB)
	default:
		log.Fatal("❌ Migration failed: unsupported database type")
	}
	//err = dbConn.Gorm.AutoMigrate(&products.Product{}, &orders.Order{})
	//if err != nil {
	//	log.Fatal("❌ Migration failed:", err)
	//}

	log.Println("✅ Database migrated successfully!")
}
