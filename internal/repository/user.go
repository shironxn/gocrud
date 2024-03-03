package repository

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(req domain.UserRequest) (*domain.User, error) {
	entity := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err := u.db.Create(entity).Error
	return &entity, err
}

func (u *UserRepository) GetAll() ([]domain.User, error) {
	var entity []domain.User
	err := u.db.Find(&entity).Error
	return entity, err
}

func (u *UserRepository) GetByID(id uint) (*domain.User, error) {
	var entity domain.User
	err := u.db.First(&entity, id).Error
	return &entity, err
}

func (u *UserRepository) Update(entity *domain.User, req domain.UserRequest) (*domain.User, error) {
	err := u.db.Model(&entity).Updates(req).Error
	return entity, err
}

func (u *UserRepository) Delete(entity *domain.User) error {
	err := u.db.Delete(entity).Error
	return err
}
