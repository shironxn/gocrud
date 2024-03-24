package route

import (
	"gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type WelcomeRoute struct {
}

func NewWelcomeRoute() WelcomeRoute {
	return WelcomeRoute{}
}

func (a *WelcomeRoute) Route(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
			Message: "welcome to gocrud by shironxn",
		})
	})
}
