package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"not null;uniqueIndex"`
	Email     string `gorm:"not null;uniqueIndex"`
	Bio       string
	AvatarURL string
	Password  string `gorm:"not null"`
	Notes     []Note `gorm:"not null"`
}

type UserDetails struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRequest struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" validate:"omitempty,min=4,max=30"`
	Email     string `json:"email" validate:"omitempty,email"`
	Bio       string `json:"bio" validate:"max=50"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url,image"`
	Password  string `json:"password" validate:"omitempty,min=8,max=100"`
}

type UserQuery struct {
	Name string `query:"name"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserPaginationResponse struct {
	Users    []UserResponse `json:"users"`
	Metadata Metadata       `json:"metadata"`
}
