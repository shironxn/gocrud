package port

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Auth() fiber.Handler
}
