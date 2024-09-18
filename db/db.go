package db

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase() (*Database, error) {
	database, err := sqlx.Open("pgx", os.Getenv("POSTGRESQL_OPEN_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	return &Database{
		db: database,
	}, nil
}

// Close a wrapper
func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
