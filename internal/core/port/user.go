package port

import (
	"gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	Create(req domain.UserRegisterRequest) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(req domain.UserRequest, entity *domain.User) (*domain.User, error)
	Delete(entity *domain.User) error
}

type UserService interface {
	Create(req domain.UserRegisterRequest) (*domain.User, error)
	Login(req domain.UserLoginRequest) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(req domain.UserRequest) (*domain.User, error)
	Update(req domain.UserRequest, claims domain.Claims) (*domain.User, error)
	Delete(req domain.UserRequest, claims domain.Claims) error
}

type UserHandler interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	GetCurrent(ctx *fiber.Ctx) error
	GetAll(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
