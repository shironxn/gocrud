package repository

import (
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NoteRepository struct {
	db         *gorm.DB
	pagination util.Pagination
}

func NewNoteRepository(db *gorm.DB, pagination util.Pagination) port.NoteRepository {
	return &NoteRepository{
		db:         db,
		pagination: pagination,
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

func (n *NoteRepository) GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error) {
	var entity []domain.Note
	if err := n.db.
		Model(entity).
		Where(&req).
		Count(&metadata.TotalRecords).
		Scopes(n.pagination.Paginate(metadata)).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (n *NoteRepository) GetByID(req domain.NoteRequest) (*domain.Note, error) {
	var entity *domain.Note
	return entity, n.db.First(&entity, req.ID).Error
}

func (n *NoteRepository) Update(req domain.NoteRequest, entity *domain.Note) (*domain.Note, error) {
	return entity, n.db.Model(&entity).Updates(req).Error
}

func (n *NoteRepository) Delete(entity *domain.Note) error {
	return n.db.Delete(&entity).Error
}
