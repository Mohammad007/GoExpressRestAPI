package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	db *sql.DB
}

func (p *PostgreSQL) Connect() error {
	db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *PostgreSQL) Close() error {
	return p.db.Close()
}

func (p *PostgreSQL) Ping() error {
	return p.db.Ping()
}