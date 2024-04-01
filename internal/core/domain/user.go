package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null;uniqueIndex"`
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null"`
	Notes    []Note `gorm:"not null"`
}

type UserDetails struct {
	Token     string
	ExpiredAt string
}

type UserRequest struct {
	ID       uint
	Name     string `validate:"required,min=4,max=30"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=100"`
}

type UserRegisterRequest struct {
	Name     string `validate:"required,min=4,max=30"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=100"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
