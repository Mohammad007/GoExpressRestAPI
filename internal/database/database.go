package database

import (
    "context"
    "fmt"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
)

type Database interface {
    Connect() error
    Close() error
    CreateUser(ctx context.Context, user *models.User) error
    GetUserByID(ctx context.Context, id uint) (*models.User, error)
    GetAllUsers(ctx context.Context) ([]models.User, error)
    UpdateUser(ctx context.Context, user *models.User) error
    DeleteUser(ctx context.Context, id uint) error
}

type Config struct {
    Type     string
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    FilePath string
}

func NewDatabase(config Config) (Database, error) {
    switch config.Type {
    case "mysql":
        return NewMySQL(config)
    case "postgres":
        return NewPostgres(config)
    case "sqlite":
        return NewSQLite(config)
    case "mongodb":
        return NewMongoDB(config)
    default:
        return nil, fmt.Errorf("unsupported database type: %s", config.Type)
    }
}