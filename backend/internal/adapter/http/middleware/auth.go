package middleware

import (
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwt util.JWT
}

func NewAuthMiddleware(jwt util.JWT) port.Middleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (a *AuthMiddleware) Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cookie := ctx.Cookies("token")

		claims, err := a.jwt.ValidateToken(cookie)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized access")
		}

		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}
