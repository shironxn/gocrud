package middleware

import (
	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwt util.JWT
	cfg *config.Config
}

func NewAuthMiddleware(jwt util.JWT, cfg *config.Config) port.Middleware {
	return &AuthMiddleware{
		jwt: jwt,
		cfg: cfg,
	}
}

func (a *AuthMiddleware) Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cookie := ctx.Cookies("access-token")

		claims, err := a.jwt.ValidateToken(cookie, a.cfg.JWT.Access)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized access")
		}

		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}
