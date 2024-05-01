package repository

import (
	"errors"
	"reflect"

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

func (r *NoteRepository) Create(req domain.NoteRequest) (*domain.Note, error) {
	var entity domain.Note

	if err := r.db.Table("notes").Select("user, title").Where("user_id = ? AND title = ?", req.UserID, req.Title).Scan(&entity).Error; err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(entity, domain.Note{}) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "note with the same title already exists")
	}

	entity = domain.Note{
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		Content:     req.Content,
		Visibility:  domain.Visibility(req.Visibility),
		UserID:      req.UserID,
	}

	if err := r.db.Create(&entity).Preload("Author").Find(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *NoteRepository) GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error) {
	var entity []domain.Note

	if err := r.db.
		Model(&domain.Note{}).
		Preload("Author").
		Where(&req).
		Count(&metadata.TotalRecords).
		Scopes(r.pagination.Paginate(metadata)).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *NoteRepository) GetByID(id uint) (*domain.Note, error) {
	var entity domain.Note

	if err := r.db.Preload("Author").First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "note not found")
		}
		return nil, err
	}

	return &entity, nil
}

func (r *NoteRepository) Update(req domain.NoteUpdateRequest, note *domain.Note) (*domain.Note, error) {
	var entity domain.Note

	if err := r.db.Table("notes").Select("id, title, user_id").Where("id != ? AND title = ? AND user_id = ?", req.ID, req.Title, req.UserID).Scan(&entity).Error; err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(entity, domain.Note{}) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "note with the same title already exists")
	} else {
		entity = *note
	}

	if err := r.db.Model(&entity).Where("id = ? AND user_id = ?", req.ID, req.UserID).Updates(req).Preload("Author").Find(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "note not found")
		}
		return nil, err
	}

	return &entity, nil
}

func (r *NoteRepository) Delete(note *domain.Note) error {
	entity := note

	if err := r.db.Delete(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "note not found")
		}
		return err
	}

	return nil
}
