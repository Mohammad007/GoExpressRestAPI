package database

import (
    "context"
    "fmt"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type MySQL struct {
    db *gorm.DB
}

func NewMySQL(config Config) (*MySQL, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.User, config.Password, config.Host, config.Port, config.DBName)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &MySQL{db: db}, nil
}

func (m *MySQL) Connect() error {
    return m.db.AutoMigrate(&models.User{})
}

func (m *MySQL) Close() error {
    sqlDB, err := m.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (m *MySQL) CreateUser(ctx context.Context, user *models.User) error {
    return m.db.WithContext(ctx).Create(user).Error
}

func (m *MySQL) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := m.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}

func (m *MySQL) GetAllUsers(ctx context.Context) ([]models.User, error) {
    var users []models.User
    err := m.db.WithContext(ctx).Find(&users).Error
    return users, err
}

func (m *MySQL) UpdateUser(ctx context.Context, user *models.User) error {
    return m.db.WithContext(ctx).Save(user).Error
}

func (m *MySQL) DeleteUser(ctx context.Context, id uint) error {
    return m.db.WithContext(ctx).Delete(&models.User{}, id).Error
}