package repository

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
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

func (u *UserRepository) Create(req domain.UserRegisterRequest) (*domain.User, error) {
	entity := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := u.db.Create(&entity).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			if strings.Contains(mysqlErr.Message, "name") {
				return nil, fiber.NewError(fiber.StatusBadRequest, "user with the same name already exists")
			}
			if strings.Contains(mysqlErr.Message, "email") {
				return nil, fiber.NewError(fiber.StatusBadRequest, "user with the same email already exists")
			}
		}
		return nil, err
	}
	return &entity, nil
}

func (u *UserRepository) GetAll(req domain.UserQuery, metadata *domain.Metadata) ([]domain.User, error) {
	var entities []domain.User
	if err := u.db.Model(&domain.User{}).
		Where(&req).
		Count(&metadata.TotalRecords).
		Scopes(u.pagination.Paginate(metadata)).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (u *UserRepository) GetByID(req domain.UserRequest) (*domain.User, error) {
	var entity domain.User
	if err := u.db.First(&entity, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, err
	}
	return &entity, nil
}

func (u *UserRepository) GetByEmail(req domain.UserRequest) (*domain.User, error) {
	var entity domain.User
	if err := u.db.Where("email = ?", req.Email).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, err
	}
	return &entity, nil
}

func (u *UserRepository) Update(req domain.UserRequest, entity *domain.User) (*domain.User, error) {
	if err := u.db.Model(entity).Updates(req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, err
	}
	return entity, nil
}

func (u *UserRepository) Delete(entity *domain.User) error {
	if err := u.db.Delete(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return err
	}
	return nil
}
