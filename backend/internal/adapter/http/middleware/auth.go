package middleware

import (
	"time"

	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	service port.AuthService
	jwt     util.JWT
	cfg     *config.Config
}

func NewAuthMiddleware(service port.AuthService, jwt util.JWT, cfg *config.Config) port.Middleware {
	return &AuthMiddleware{
		service: service,
		jwt:     jwt,
		cfg:     cfg,
	}
}

func (m *AuthMiddleware) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access-token")
		refreshToken := c.Cookies("refresh-token")

		claims, err := m.jwt.ValidateToken(accessToken, m.cfg.JWT.Access)
		if err != nil {
			if refreshToken == "" {
				return fiber.NewError(fiber.StatusUnauthorized, "unauthorized access")
			}

			newAccessToken, claims, err := m.service.Refresh(refreshToken)
			if err != nil {
				return fiber.NewError(fiber.StatusUnauthorized, "unauthorized access")
			}

			c.Cookie(&fiber.Cookie{
				Name:     "access-token",
				Value:    *newAccessToken,
				Path:     "/",
				HTTPOnly: true,
				Expires:  time.Now().Add(10 * time.Minute),
				SameSite: func(dev string) string {
					if dev == "true" {
						return fiber.CookieSameSiteLaxMode
					}
					return fiber.CookieSameSiteNoneMode
				}(m.cfg.Server.Dev),
			})
			c.Locals("claims", claims)

			return c.Next()
		}

		c.Locals("claims", claims)

		return c.Next()
	}
}
