package repository

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"

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
	entity := domain.Note{
		Title:      req.Title,
		Content:    req.Content,
		Visibility: req.Visibility,
	}
	return &entity, n.db.Create(&entity).Error
}

func (n *NoteRepository) GetAll() ([]domain.Note, error) {
	var entity []domain.Note
	return entity, n.db.Find(&entity).Error
}

func (n *NoteRepository) GetByID(id uint) (*domain.Note, error) {
	var entity *domain.Note
	return entity, n.db.First(&entity, id).Error
}

func (n *NoteRepository) Update(req domain.NoteRequest, entity *domain.Note) (*domain.Note, error) {
	return entity, n.db.Model(&req).Updates(req).Error
}

func (n *NoteRepository) Delete(entity *domain.Note) error {
	return n.db.Delete(&entity).Error
}
