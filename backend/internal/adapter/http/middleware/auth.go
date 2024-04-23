package middleware

import (
	"time"

	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
			log.Info(newAccessToken)
			c.Cookie(&fiber.Cookie{
				Name:     "access-token",
				Value:    *newAccessToken,
				Path:     "/",
				HTTPOnly: true,
				Expires:  time.Now().Add(10 * time.Minute),
				SameSite: func(mode string) string {
					if mode != "DEV" {
						return fiber.CookieSameSiteNoneMode
					}
					return fiber.CookieSameSiteLaxMode
				}(m.cfg.Server.Mode),
			})
			c.Locals("claims", claims)

			return c.Next()
		}

		c.Locals("claims", claims)

		return c.Next()
	}
}
