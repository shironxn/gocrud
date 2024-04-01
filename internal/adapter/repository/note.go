package repository

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) port.NoteRepository {
	return &NoteRepository{
		db: db,
	}
}

func (n *NoteRepository) Create(req domain.NoteRequest) (*domain.Note, error) {
	var user domain.User
	err := n.db.Preload("Notes").Find(&user, req.UserID).Error
	if err != nil {
		return nil, err
	}
	for _, note := range user.Notes {
		if note.Title == req.Title {
			return nil, fiber.NewError(fiber.StatusBadRequest, gorm.ErrDuplicatedKey.Error())
		}
	}
	entity := domain.Note{
		Title:      req.Title,
		Content:    req.Content,
		Visibility: req.Visibility,
		UserID:     req.UserID,
	}
	return &entity, n.db.Create(&entity).Error
}

func (n *NoteRepository) GetAll() ([]domain.Note, error) {
	var entity []domain.Note
	return entity, n.db.Find(&entity, "visibility = ?", "public").Error
}

func (n *NoteRepository) GetByID(id uint) (*domain.Note, error) {
	var entity *domain.Note
	return entity, n.db.First(&entity, id).Error
}

func (n *NoteRepository) Update(req domain.NoteRequest, entity *domain.Note) (*domain.Note, error) {
	return entity, n.db.Model(&entity).Updates(req).Error
}

func (n *NoteRepository) Delete(entity *domain.Note) error {
	return n.db.Delete(&entity).Error
}
