package route

import (
	"github.com/shironxn/gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute struct {
	handler    port.UserHandler
	middleware port.Middleware
}

func NewAuthRoute(handler port.UserHandler, middleware port.Middleware) AuthRoute {
	return AuthRoute{
		handler:    handler,
		middleware: middleware,
	}
}

func (a *AuthRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/auth")
	v1.Post("/register", a.handler.Register)
	v1.Post("/login", a.handler.Login)
	v1.Post("/logout", a.middleware.Auth(), a.handler.Logout)
}
