package route

import (
	"github.com/shironxn/gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type UserRoute struct {
	handler    port.UserHandler
	middleware port.Middleware
}

func NewUserRoute(handler port.UserHandler, middleware port.Middleware) UserRoute {
	return UserRoute{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *UserRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/users")
	v1.Get("/", r.handler.GetAll)
	v1.Get("/me", r.middleware.Auth(), r.handler.GetMe)
	v1.Get("/:id", r.handler.GetByID)
	v1.Put("/:id", r.middleware.Auth(), r.handler.Update)
	v1.Delete(":id", r.middleware.Auth(), r.handler.Delete)
}
