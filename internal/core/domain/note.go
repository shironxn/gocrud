package domain

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	gorm.Model `gorm:"not null"`
	Title      string `gorm:"not null;unique"`
	Content    string `gorm:"not null"`
	Visibility string `gorm:"not null;type:enum('public','private');default:'private'"`
	UserID     uint   `gorm:"not null"`
	User       User   `gorm:"foreignKey:UserID;references:ID"`
}

type NoteRequest struct {
	ID         uint
	Title      string `validate:"required"`
	Content    string `validate:"required"`
	Visibility string `validate:"oneof=private public"`
	UserID     uint   `validate:"required"`
}

type NoteResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Visibility string    `json:"visibility"`
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
