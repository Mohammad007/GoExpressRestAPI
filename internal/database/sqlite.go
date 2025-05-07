package database

import (
    "context"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type SQLite struct {
    db *gorm.DB
}

func NewSQLite(config Config) (*SQLite, error) {
    db, err := gorm.Open(sqlite.Open(config.FilePath), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &SQLite{db: db}, nil
}

func (s *SQLite) Connect() error {
    return s.db.AutoMigrate(&models.User{})
}

func (s *SQLite) Close() error {
    sqlDB, err := s.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (s *SQLite) CreateUser(ctx context.Context, user *models.User) error {
    return s.db.WithContext(ctx).Create(user).Error
}

func (s *SQLite) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := s.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}

func (s *SQLite) GetAllUsers(ctx context.Context) ([]models.User, error) {
    var users []models.User
    err := s.db.WithContext(ctx).Find(&users).Error
    return users, err
}

func (s *SQLite) UpdateUser(ctx context.Context, user *models.User) error {
    return s.db.WithContext(ctx).Save(user).Error
}

func (s *SQLite) DeleteUser(ctx context.Context, id uint) error {
    return s.db.WithContext(ctx).Delete(&models.User{}, id).Error
}