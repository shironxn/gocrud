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

func (h *UserRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/users")
	v1.Get("/", h.handler.GetAll)
	v1.Get("/:id", h.handler.GetByID)
	v1.Put("/:id", h.middleware.Auth(), h.handler.Update)
	v1.Delete(":id", h.middleware.Auth(), h.handler.Delete)
}
