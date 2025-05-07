package database

import (
    "context"
    "fmt"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Postgres struct {
    db *gorm.DB
}

func NewPostgres(config Config) (*Postgres, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.Host, config.Port, config.User, config.Password, config.DBName)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Postgres{db: db}, nil
}

func (p *Postgres) Connect() error {
    return p.db.AutoMigrate(&models.User{})
}

func (p *Postgres) Close() error {
    sqlDB, err := p.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (p *Postgres) CreateUser(ctx context.Context, user *models.User) error {
    return p.db.WithContext(ctx).Create(user).Error
}

func (p *Postgres) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := p.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}

func (p *Postgres) GetAllUsers(ctx context.Context) ([]models.User, error) {
    var users []models.User
    err := p.db.WithContext(ctx).Find(&users).Error
    return users, err
}

func (p *Postgres) UpdateUser(ctx context.Context, user *models.User) error {
    return p.db.WithContext(ctx).Save(user).Error
}

func (p *Postgres) DeleteUser(ctx context.Context, id uint) error {
    return p.db.WithContext(ctx).Delete(&models.User{}, id).Error
}