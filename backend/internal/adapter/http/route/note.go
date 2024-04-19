package route

import (
	"github.com/shironxn/gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type NoteRoute struct {
	handler    port.NoteHandler
	middleware port.Middleware
}

func NewNoteRoute(handler port.NoteHandler, middleware port.Middleware) NoteRoute {
	return NoteRoute{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *NoteRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/notes")
	v1.Post("/", r.middleware.Auth(), r.handler.Create)
	v1.Get("/", r.handler.GetAll)
	v1.Get("/:id", r.handler.GetByID)
	v1.Put("/:id", r.middleware.Auth(), r.handler.Update)
	v1.Delete("/:id", r.middleware.Auth(), r.handler.Delete)
}
