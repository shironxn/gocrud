package route

import (
	"gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "gocrud/docs"

	"github.com/gofiber/swagger"
)

type WelcomeRoute struct {
}

func NewWelcomeRoute() WelcomeRoute {
	return WelcomeRoute{}
}

func (a *WelcomeRoute) Route(app *fiber.App) {
	app.Use(logger.New())
	app.Get("/api/v1/docs/*", swagger.HandlerDefault)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
			Message: "welcome to gocrud by shironxn",
		})
	})
}
