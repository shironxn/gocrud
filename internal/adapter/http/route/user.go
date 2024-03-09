package route

import (
	"gocrud/internal/adapter/http/middleware"
	"gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type UserRoute struct {
	userHandler port.UserHandler
	middleware  *middleware.AuthMiddleware
}

func NewUserRoute(userHandler port.UserHandler, middleware *middleware.AuthMiddleware) *UserRoute {
	return &UserRoute{
		userHandler: userHandler,
		middleware:  middleware,
	}
}

func (u *UserRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/user", u.middleware.Auth())
	v1.Get("/current", u.userHandler.GetCurrent)
	v1.Get("/", u.userHandler.GetAll)
	v1.Get("/:id", u.userHandler.GetByID)
	v1.Post("/", u.userHandler.Register)
	v1.Put("/:id", u.userHandler.Update)
	v1.Delete(":id", u.userHandler.Delete)
}
