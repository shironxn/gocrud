package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shironxn/gocrud/internal/config"

	_ "github.com/shironxn/gocrud/docs"

	"github.com/gofiber/swagger"
)

type InitRoute struct {
	cfg *config.Config
}

func NewInitRoute(cfg *config.Config) InitRoute {
	return InitRoute{
		cfg: cfg,
	}
}

func (r *InitRoute) Route(app *fiber.App) {
	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     r.cfg.Server.Web,
			AllowCredentials: true,
		},
	))
	app.Use(logger.New())

	app.Get("/api/v1/docs/*", swagger.HandlerDefault)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON("Welcome to gocrud by shironxn")
	})
}
