package db

import "database/sql"

// Database defines common database operations
type Database interface {
	GetConnection() (*sql.DB, error)
}
