package repository

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) port.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Register(req domain.AuthRegisterRequest) (*domain.User, error) {
	entity := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := r.db.Create(&entity).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "idx_users_name":
				return nil, fiber.NewError(fiber.StatusBadRequest, "user with the same name already exists")
			case "idx_users_email":
				return nil, fiber.NewError(fiber.StatusBadRequest, "user with the same email already exists")
			}
		}
		return nil, err
	}
	return &entity, nil
}

func (r *AuthRepository) GetByEmail(email string) (*domain.User, error) {
	var entity domain.User
	if err := r.db.Where("email = ?", email).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, err
	}
	return &entity, nil
}

func (r *AuthRepository) GetRefreshToken(id uint) (*domain.RefreshToken, error) {
	var entity domain.RefreshToken
	if err := r.db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "refresh token not found")
		}
		return nil, err
	}
	return &entity, nil

}

func (r *AuthRepository) StoreRefreshToken(id uint, token string) error {
	var entity = domain.RefreshToken{
		UserID: id,
		Token:  token,
	}
	return r.db.Save(&entity).Error
}

func (r *AuthRepository) DeleteRefreshToken(entity domain.RefreshToken) error {
	return r.db.Delete(&entity).Error
}
