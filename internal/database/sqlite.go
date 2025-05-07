package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func (s *SQLite) Connect() error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) Ping() error {
	return s.db.Ping()
}