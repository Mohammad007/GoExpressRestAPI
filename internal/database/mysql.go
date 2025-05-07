package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	db *sql.DB
}

func (m *MySQL) Connect() error {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *MySQL) Close() error {
	return m.db.Close()
}

func (m *MySQL) Ping() error {
	return m.db.Ping()
}