package repository

import (
	"errors"

	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	db         *gorm.DB
	pagination util.Pagination
}

func NewUserRepository(db *gorm.DB, pagination util.Pagination) port.UserRepository {
	return &UserRepository{
		db:         db,
		pagination: pagination,
	}
}

func (r *UserRepository) GetAll(req domain.UserQuery, metadata *domain.Metadata) ([]domain.User, error) {
	var entities []domain.User
	if err := r.db.Model(&domain.User{}).
		Where(&req).
		Count(&metadata.TotalRecords).
		Scopes(r.pagination.Paginate(metadata)).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	var entity domain.User
	if err := r.db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, err
	}
	return &entity, nil
}

func (r *UserRepository) Update(req domain.UserRequest, entity *domain.User) (*domain.User, error) {
	if err := r.db.Model(entity).Updates(req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, err
	}
	return entity, nil
}

func (r *UserRepository) Delete(entity *domain.User) error {
	if err := r.db.Delete(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return err
	}
	return nil
}
