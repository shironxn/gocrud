package port

import (
	"github.com/shironxn/gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	GetAll(req domain.UserQuery, metdata *domain.Metadata) ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(req domain.UserRequest, entity *domain.User) (*domain.User, error)
	Delete(entity *domain.User) error
}

type UserService interface {
	GetAll(req domain.UserQuery, metdata *domain.Metadata) ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(req domain.UserRequest, claims domain.Claims) (*domain.User, error)
	Delete(req domain.UserRequest, claims domain.Claims) error
}

type UserHandler interface {
	GetAll(ctx *fiber.Ctx) error
	GetMe(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
