package port

import (
	"gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	Create(req domain.UserRequest) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(entity *domain.User, req domain.UserRequest) (*domain.User, error)
	Delete(entity *domain.User) error
}

type UserService interface {
	Create(req domain.UserRequest) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(req domain.UserRequest) (*domain.User, error)
	Update(req domain.UserRequest, claims domain.Claims) (*domain.User, error)
	Delete(req domain.UserRequest, claims domain.Claims) error
}

type UserHandler interface {
	Create(ctx *fiber.Ctx) error
	GetAll(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
