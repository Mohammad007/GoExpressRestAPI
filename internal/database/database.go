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
        // Try to use the regular SQLite implementation first
        db, err := NewSQLite(config)
        if err != nil {
            // If we get a CGO error, try the pure Go implementation
            if err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" {
                // Inform user we're falling back to pure Go implementation
                fmt.Println("WARNING: CGO is disabled. Falling back to pure Go SQLite implementation.")
                fmt.Println("To use standard SQLite, rebuild with CGO_ENABLED=1")
                
                // Try the pure Go implementation
                pureDb, pureErr := NewSQLitePure(config)
                if pureErr != nil {
                    // If pure implementation also fails, provide comprehensive error
                    return nil, fmt.Errorf("database error: SQLite requires CGO_ENABLED=1 to work.\n" +
                        "Pure Go SQLite alternative also failed: %v\n" +
                        "Please either:\n" +
                        "1. Rebuild with CGO_ENABLED=1\n" +
                        "2. Configure an alternative database like MySQL or PostgreSQL\n" +
                        "3. Install the pure Go SQLite dependency with 'go get github.com/glebarez/sqlite'", pureErr)
                }
                return pureDb, nil
            }
            return nil, err
        }
        return db, nil
    case "sqlite-pure":
        // Explicitly use the pure Go implementation
        db, err := NewSQLitePure(config)
        if err != nil {
            // Provide more helpful error for pure SQLite implementation
            return nil, fmt.Errorf("failed to initialize pure Go SQLite: %v\n" +
                "Make sure you have installed the required dependency with:\n" +
                "go get github.com/glebarez/sqlite\ngo mod tidy", err)
        }
        return db, nil
    case "mongodb":
        return NewMongoDB(config)
    default:
        return nil, fmt.Errorf("unsupported database type: %s\nSupported types: sqlite, sqlite-pure, mysql, postgres, mongodb", config.Type)
    }
}