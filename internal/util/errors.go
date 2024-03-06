package util

import (
	"gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
}

func (r *Response) Success(ctx *fiber.Ctx, res domain.SuccessResponse) error {
	response := domain.SuccessResponse{
		Message: res.Message,
		Data:    res.Data,
	}

	return ctx.Status(res.Status).JSON(response)
}

func (r *Response) Error(ctx *fiber.Ctx, res domain.ErrorResponse) error {
	response := domain.ErrorResponse{
		Message: res.Message,
		Errors:  res.Errors,
	}

	return ctx.Status(res.Status).JSON(response)
}
