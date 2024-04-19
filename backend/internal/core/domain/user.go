package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null;uniqueIndex"`
	Email        string `gorm:"not null;uniqueIndex"`
	Bio          string
	AvatarURL    string
	Password     string `gorm:"not null"`
	RefreshToken RefreshToken
	Notes        []Note `gorm:"not null"`
}

type RefreshToken struct {
	UserID uint `gorm:"primarykey"`
	Token  string
}

type UserToken struct {
	AccessToken  string  `json:"access_token,omitempty"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	Claims       *Claims `json:"claims,omitempty"`
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
	Name    string `query:"name"`
	Details bool   `query:"details"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Bio       string     `json:"bio,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserToken *UserToken `json:"tokens,omitempty"`
}

type UserPaginationResponse struct {
	Users    []UserResponse `json:"users"`
	Metadata Metadata       `json:"metadata"`
}
