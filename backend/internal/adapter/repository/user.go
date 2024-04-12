package repository

import (
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

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

func (u *UserRepository) Create(req domain.UserRegisterRequest) (*domain.User, error) {
	entity := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	return &entity, u.db.Create(&entity).Error
}

func (u *UserRepository) GetAll(req domain.UserQuery, metadata *domain.Metadata) ([]domain.User, error) {
	var entity []domain.User
	if err := u.db.Model(entity).
		Where(&req).
		Count(&metadata.TotalRecords).
		Scopes(u.pagination.Paginate(metadata)).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *UserRepository) GetByID(req domain.UserRequest) (*domain.User, error) {
	var entity domain.User
	return &entity, u.db.First(&entity, req.ID).Error
}

func (u *UserRepository) GetByEmail(req domain.UserRequest) (*domain.User, error) {
	var entity domain.User
	return &entity, u.db.Where("email = ?", req.Email).First(&entity).Error
}

func (u *UserRepository) Update(req domain.UserRequest, entity *domain.User) (*domain.User, error) {
	return entity, u.db.Model(&entity).Updates(req).Error
}

func (u *UserRepository) Delete(entity *domain.User) error {
	return u.db.Select("Notes").Delete(&entity).Error
}
