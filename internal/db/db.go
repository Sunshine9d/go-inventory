package db

import "database/sql"

type Database interface {
	NewConnection() (*sql.DB, error)
}

type DBSql struct {
	DB *sql.DB
}
