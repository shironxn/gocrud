package domain

import (
	"time"

	"gorm.io/gorm"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type Note struct {
	gorm.Model
	Title       string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	CoverURL    string     `gorm:"not null"`
	Content     string     `gorm:"not null"`
	Visibility  Visibility `gorm:"not null;default:'private'" sql:"type:visibility"`
	UserID      uint       `gorm:"not null"`
	Author      User       `gorm:"foreignKey:UserID"`
}

type NoteRequest struct {
	ID          uint   `json:"id"`
	Title       string `json:"title" validate:"required,max=25" conform:"name,title,alpha"`
	Description string `json:"description" validate:"required,max=50" conform:"trim"`
	CoverURL    string `json:"cover_url" validate:"required,url,image" conform:"trim"`
	Content     string `json:"content" validate:"required" conform:"trim"`
	Visibility  string `json:"visibility" validate:"required,oneof=private public"`
	UserID      uint   `json:"user_id"`
}

type NoteUpdateRequest struct {
	ID          uint   `json:"id"`
	Title       string `json:"title" validate:"omitempty,max=25" conform:"name,title"`
	Description string `json:"description" validate:"omitempty,max=50" conform:"trim"`
	CoverURL    string `json:"cover_url" validate:"omitempty,url,image" conform:"trim"`
	Content     string `json:"content" validate:"omitempty" conform:"trim"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=private public"`
	UserID      uint   `json:"user_id"`
}

type NoteQuery struct {
	Title      string `query:"title"`
	Visibility string `query:"visibility"`
	UserID     int    `query:"user_id"`
}

type NoteAuthor struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type NoteResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CoverURL    string     `json:"cover_url"`
	Content     string     `json:"content"`
	Visibility  string     `json:"visibility"`
	Author      NoteAuthor `json:"author"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type NotePaginationResponse struct {
	Notes    []NoteResponse `json:"notes"`
	Metadata Metadata       `json:"metadata"`
}
