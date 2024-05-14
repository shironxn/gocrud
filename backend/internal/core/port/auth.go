package port

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shironxn/blanknotes/internal/core/domain"
)

type AuthRepository interface {
	Register(req domain.AuthRegisterRequest) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetRefreshToken(userID uint) (*domain.RefreshToken, error)
	StoreRefreshToken(userID uint, token string) error
	DeleteRefreshToken(entity domain.RefreshToken) error
}

type AuthService interface {
	Register(req domain.AuthRegisterRequest) (*domain.User, error)
	Login(req domain.AuthLoginRequest) (*domain.User, *domain.UserToken, error)
	Logout(userID uint) error
	Refresh(token string) (*string, *domain.Claims, error)
}

type AuthHandler interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
}
