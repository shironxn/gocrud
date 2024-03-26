package route

import (
	"gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
	_ "github.com/shironxn/gocrud/docs"
)

type WelcomeRoute struct {
}

func NewWelcomeRoute() WelcomeRoute {
	return WelcomeRoute{}
}

func (a *WelcomeRoute) Route(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
			Message: "welcome to gocrud by shironxn",
		})
	})

}
