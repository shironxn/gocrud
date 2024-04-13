package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shironxn/gocrud/internal/core/domain"

	_ "github.com/shironxn/gocrud/docs"

	"github.com/gofiber/swagger"
)

type InitRoute struct {
}

func NewInitRoute() InitRoute {
	return InitRoute{}
}

func (a *InitRoute) Route(app *fiber.App) {
	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     "http://localhost:3001",
			AllowCredentials: true,
		},
	))
	app.Use(logger.New())

	app.Get("/api/v1/docs/*", swagger.HandlerDefault)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
			Message: "welcome to gocrud by shironxn",
		})
	})
}
