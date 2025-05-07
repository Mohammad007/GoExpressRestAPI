package database

import (
	"context"
	"fmt"
	"github.com/Mohammad007/GoExpressRestAPI/internal/models"
	"github.com/glebarez/sqlite" // Pure Go SQLite implementation
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SQLitePure is a fallback implementation when SQLite with CGO is not available
// It can use either a pure Go SQLite implementation or fall back to another database type
type SQLitePure struct {
	db *gorm.DB
}

// NewSQLitePure creates a fallback database connection when SQLite with CGO is not available
func NewSQLitePure(config Config) (*SQLitePure, error) {
	// First try to use the pure Go SQLite implementation
	if config.FilePath != "" {
		// Use a defer/recover to catch panics that might occur if the package is not installed
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("ERROR: Failed to initialize pure Go SQLite implementation.")
					fmt.Println("Please run 'go get github.com/glebarez/sqlite' and 'go mod tidy' to install the required dependency.")
				}
			}()
			
			fmt.Println("INFO: Using pure Go SQLite implementation (github.com/glebarez/sqlite)")
			db, err := gorm.Open(sqlite.Open(config.FilePath), &gorm.Config{})
			if err == nil {
				// Only set the db if no error occurred
				return &SQLitePure{db: db}, nil
			}
			fmt.Println("WARNING: Failed to initialize pure Go SQLite, falling back to alternative database:", err)
		}()
	}

	// If pure Go SQLite fails or isn't configured, try MySQL as fallback
	if config.Host != "" && config.User != "" {
		fmt.Println("WARNING: CGO is disabled. SQLite requires CGO to work.")
		fmt.Println("WARNING: Attempting to use MySQL as fallback. To use SQLite, rebuild with CGO_ENABLED=1")
		
		// Use MySQL if configuration is available
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.DBName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL fallback: %w\nOriginal error: SQLite requires CGO", err)
		}
		return &SQLitePure{db: db}, nil
	}
	
	// If no fallback config is available, return a clear error message
	return nil, fmt.Errorf("database error: SQLite requires CGO_ENABLED=1 to work.\n" +
		"Please either:\n" +
		"1. Rebuild with CGO_ENABLED=1\n" +
		"2. Configure an alternative database like MySQL or PostgreSQL\n" +
		"3. Use 'sqlite-pure' with a pure Go SQLite alternative (requires additional setup)")
}
}

func (s *SQLitePure) Connect() error {
	return s.db.AutoMigrate(&models.User{})
}

func (s *SQLitePure) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *SQLitePure) CreateUser(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *SQLitePure) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (s *SQLitePure) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := s.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (s *SQLitePure) UpdateUser(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *SQLitePure) DeleteUser(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&models.User{}, id).Error
}