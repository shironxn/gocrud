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

func (n *NoteRoute) Route(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1/notes")
	v1.Post("/", n.middleware.Auth(), n.handler.Create)
	v1.Get("/", n.handler.GetAll)
	v1.Get("/:id", n.handler.GetByID)
	v1.Put("/:id", n.middleware.Auth(), n.handler.Update)
	v1.Delete("/:id", n.middleware.Auth(), n.handler.Delete)
}
