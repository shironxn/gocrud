package route

import (
	"gocrud/internal/core/port"

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

func (u *UserRoute) Route(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1/user", u.middleware.Auth())

	v1.Get("/current", u.handler.GetCurrent)
	v1.Get("/", u.handler.GetAll)
	v1.Get("/:id", u.handler.GetByID)
	v1.Put("/:id", u.handler.Update)
	v1.Delete(":id", u.handler.Delete)
}
