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

func (h *NoteRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/notes")
	v1.Post("/", h.middleware.Auth(), h.handler.Create)
	v1.Get("/", h.handler.GetAll)
	v1.Get("/:id", h.handler.GetByID)
	v1.Put("/:id", h.middleware.Auth(), h.handler.Update)
	v1.Delete("/:id", h.middleware.Auth(), h.handler.Delete)
}
