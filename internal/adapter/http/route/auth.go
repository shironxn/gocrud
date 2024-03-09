package route

import (
	"gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute struct {
	userHandler port.UserHandler
}

func NewAuthRoute(userHandler port.UserHandler) *AuthRoute {
	return &AuthRoute{
		userHandler: userHandler,
	}
}

func (a *AuthRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/auth")
	v1.Post("/register", a.userHandler.Register)
	v1.Post("/login", a.userHandler.Login)
	v1.Post("/logout", a.userHandler.Logout)
}
