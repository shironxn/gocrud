package config

import (
	"github.com/shironxn/gocrud/internal/core/domain"

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

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(domain.ErrorResponse{
			Code:  code,
			Error: err.Error(),
		})
	}
}
