package mysql

import (
	"database/sql"
	"log"
)

// MigrateMySQL runs the database migrations for MySQL
func MigrateMySQL(db *sql.DB) error {
	// Create products table
	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		sku INT NOT NULL,
		price DECIMAL(10,2) NOT NULL
	);
	`
	_, err := db.Exec(productTable)
	if err != nil {
		log.Printf("❌ Failed to create products table: %v\n", err)
		return err
	}

	// Create orders table
	orderTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		customer_name VARCHAR(255) NOT NULL,
		total_price DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(orderTable)
	if err != nil {
		log.Printf("❌ Failed to create orders table: %v\n", err)
		return err
	}

	log.Println("✅ MySQL migrations completed successfully!")
	return nil
}
