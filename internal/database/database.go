package database

type Database interface {
	Connect() error
	Close() error
	Ping() error
}