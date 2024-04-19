package domain

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Content    string `gorm:"not null"`
	Visibility string `gorm:"not null;type:enum('public','private');default:'private'"`
	UserID     uint   `gorm:"not null"`
	Author     User   `gorm:"foreignKey:UserID"`
}

type NoteRequest struct {
	ID         uint   `json:"id"`
	Title      string `json:"title" validate:"max=30"`
	Content    string `json:"content"`
	Visibility string `json:"visibility" validate:"omitempty,oneof=private public"`
	UserID     uint   `json:"user_id"`
}

type NoteQuery struct {
	Title      string `query:"title"`
	Visibility string `query:"visibility"`
	UserID     string `query:"user_id"`
}

type NoteAuthor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type NoteResponse struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Visibility string     `json:"visibility"`
	Author     NoteAuthor `json:"author"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type NotePaginationResponse struct {
	Notes    []NoteResponse `json:"notes"`
	Metadata Metadata       `json:"metadata"`
}
