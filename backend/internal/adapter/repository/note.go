package repository

import (
	"errors"

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
	if err := n.db.Where("user_id = ? AND title = ?", req.UserID, req.Title).First(&domain.Note{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		return nil, fiber.NewError(fiber.StatusBadRequest, "note with the same title already exists")
	}
	entity := domain.Note{
		Title:      req.Title,
		Content:    req.Content,
		Visibility: req.Visibility,
		UserID:     req.UserID,
	}
	if err := n.db.Create(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (n *NoteRepository) GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error) {
	var entity []domain.Note
	if err := n.db.
		Model(&domain.Note{}).
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
	var entity domain.Note
	if err := n.db.First(&entity, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "note not found")
		}
		return nil, err
	}
	return &entity, nil
}

func (n *NoteRepository) Update(req domain.NoteRequest, entity *domain.Note) (*domain.Note, error) {
	if err := n.db.Where("user_id = ? AND title = ? AND id != ?", req.UserID, req.Title, entity.ID).First(&domain.Note{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		return nil, fiber.NewError(fiber.StatusBadRequest, "note with the same title already exists")
	}
	if err := n.db.Model(&entity).Updates(req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "note not found")
		}
		return nil, err
	}
	return entity, nil
}

func (n *NoteRepository) Delete(entity *domain.Note) error {
	if err := n.db.Delete(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "note not found")
		}
		return err
	}
	return nil
}
