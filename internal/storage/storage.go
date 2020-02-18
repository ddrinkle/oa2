package storage

import "database/sql"

type Crate struct {
	DB *sql.DB
}
