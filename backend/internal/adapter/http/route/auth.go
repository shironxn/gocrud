package route

import (
	"github.com/shironxn/blanknotes/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute struct {
	handler    port.AuthHandler
	middleware port.Middleware
}

func NewAuthRoute(handler port.AuthHandler, middleware port.Middleware) AuthRoute {
	return AuthRoute{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *AuthRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/auth")
	v1.Post("/register", r.handler.Register)
	v1.Post("/login", r.handler.Login)
	v1.Post("/logout", r.middleware.Auth(), r.handler.Logout)
	v1.Post("/refresh", r.handler.Refresh)
}
