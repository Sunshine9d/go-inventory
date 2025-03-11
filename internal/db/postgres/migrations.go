package postgres

import (
	"database/sql"
	"log"
)

// MigratePostgres runs the database migrations for PostgreSQL
func MigratePostgres(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity INT NOT NULL,
		price DECIMAL(10,2) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		customer_name VARCHAR(255) NOT NULL,
		total_price DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("❌ Failed to run PostgreSQL migrations: %v\n", err)
		return err
	}

	log.Println("✅ PostgreSQL migrations completed successfully!")
	return nil
}
