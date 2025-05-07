// internal/models/user.go
package models

import (
    "github.com/go-playground/validator/v10"
    "gorm.io/gorm"
    "time"
)

type UserSchema struct {
    ID        uint       `json:"id" gorm:"primaryKey" bson:"id"`
    Name      string     `json:"name" validate:"required,min=2" gorm:"type:varchar(100)" bson:"name"`
    Email     string     `json:"email" validate:"required,email" gorm:"unique;type:varchar(100)" bson:"email"`
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" bson:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" bson:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index" bson:"deleted_at"`
}

type User struct {
    UserSchema
    validator *validator.Validate
}

func NewUser(name, email string) *User {
    user := &User{
        UserSchema: UserSchema{
            Name:      name,
            Email:     email,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        validator: validator.New(),
    }
    return user
}

func (u *User) Validate() error {
    return u.validator.Struct(u)
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    return u.Validate()
}