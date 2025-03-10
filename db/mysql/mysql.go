package mysql

import "database/sql"

type MysqlDB struct {
	DB *sql.DB
}
