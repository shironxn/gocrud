package handler

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"
	"gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Response  util.Response
	Validator util.Validator
	Service   port.UserService
}

func NewUserService(validator util.Validator, service port.UserService) port.UserHandler {
	return &UserHandler{
		Validator: validator,
		Service:   service,
	}
}

func (u *UserHandler) Create(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := u.Validator.Validate(req); err != nil {
		u.Response.Error(ctx, domain.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "validation error",
			Errors:  err,
		})
		return nil
	}

	result, err := u.Service.Create(req)
	if err != nil {
		u.Response.Error(ctx, domain.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "cannot create user",
			Errors:  err,
		})
		return nil
	}

	data := domain.UserResponse{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	u.Response.Success(ctx, domain.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "successfully create user",
		Data:    data,
	})

	return nil
}

func (u *UserHandler) GetAll(ctx *fiber.Ctx) error {
	result, err := u.Service.GetAll()
	if err != nil {
		u.Response.Error(ctx, domain.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "failed to get all user data",
			Errors:  err,
		})
	}

	var data []domain.UserResponse
	for _, user := range result {
		data = append(data,
			domain.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
	}

	u.Response.Success(ctx, domain.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "successfully get all user data",
		Data:    data,
	})

	return nil
}

func (u *UserHandler) GetByID(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	err := ctx.ParamsParser(&req.ID)
	if err != nil {
		return err
	}

	result, err := u.Service.GetByID(req)
	if err != nil {
		u.Response.Error(ctx, domain.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "failed to get user by id",
			Errors:  err,
		})
	}

	data := domain.UserResponse{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	u.Response.Success(ctx, domain.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "successfully get user by id",
		Data:    data,
	})

	return nil
}

func (u *UserHandler) Update(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req.ID); err != nil {
		return err
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	var claims domain.Claims
	if err := ctx.CookieParser(&claims); err != nil {
		return err
	}

	result, err := u.Service.Update(req, claims)
	if err != nil {
		u.Response.Error(ctx, domain.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "failed to update user by id",
			Errors:  err,
		})
	}

	data := domain.UserResponse{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	u.Response.Success(ctx, domain.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "successfully update user by id",
		Data:    data,
	})

	return nil
}

func (u *UserHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req.ID); err != nil {
		return err
	}

	var claims domain.Claims
	if err := ctx.CookieParser(&claims); err != nil {
		return err
	}

	err := u.Service.Delete(req, claims)
	if err != nil {
		u.Response.Error(ctx, domain.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "failed to delete user by id",
			Errors:  err,
		})
	}

	u.Response.Success(ctx, domain.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "successfully delete user by id",
		Data:    nil,
	})

	return nil
}
