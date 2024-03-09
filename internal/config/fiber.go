package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})
}

func ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
