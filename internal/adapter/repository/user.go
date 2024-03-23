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

func (u *UserRepository) Create(req domain.UserRegisterRequest) (*domain.User, error) {
	entity := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	return &entity, u.db.Create(&entity).Error
}

func (u *UserRepository) GetAll() ([]domain.User, error) {
	var entity []domain.User
	return entity, u.db.Find(&entity).Error
}

func (u *UserRepository) GetByID(id uint) (*domain.User, error) {
	var entity domain.User
	return &entity, u.db.First(&entity, id).Error
}

func (u *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var entity domain.User
	return &entity, u.db.Where("email = ?", email).First(&entity).Error
}

func (u *UserRepository) Update(req domain.UserRequest, entity *domain.User) (*domain.User, error) {
	return entity, u.db.Model(&entity).Updates(req).Error
}

func (u *UserRepository) Delete(entity *domain.User) error {
	return u.db.Delete(&entity).Error
}
