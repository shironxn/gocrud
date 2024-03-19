package handler

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"
	"gocrud/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service   port.UserService
	validator *util.Validator
	jwt       *util.JWT
}

func NewUserHandler(service port.UserService, validator *util.Validator, jwt *util.JWT) port.UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator,
		jwt:       jwt,
	}
}

func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	var req domain.UserRegisterRequest

	if cookie := ctx.Cookies("token"); cookie != "" {
		return fiber.NewError(fiber.StatusBadRequest, "user already registered")
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := u.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := u.service.Create(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully create user",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	var req domain.UserLoginRequest

	if cookie := ctx.Cookies("token"); cookie != "" {
		return fiber.NewError(fiber.StatusBadRequest, "user already logged in")
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := u.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := u.service.Login(req)
	if err != nil {
		return err
	}

	_, err = u.jwt.GenerateToken(ctx, result)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully login user",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}})
}

func (u *UserHandler) Logout(ctx *fiber.Ctx) error {
	if cookie := ctx.Cookies("token"); cookie == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user already logged out")
	}

	cookieExpire := time.Now().Add(-time.Hour * 24)
	ctx.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: cookieExpire,
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully logout user",
	})
}

func (u *UserHandler) GetCurrent(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	claims := ctx.Locals("claims").(*domain.Claims)
	req.ID = claims.UserID

	result, err := u.service.GetByID(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully get current user data",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

func (u *UserHandler) GetAll(ctx *fiber.Ctx) error {
	result, err := u.service.GetAll()
	if err != nil {
		return err
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

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully get all user data",
		Data:    data,
	})
}

func (u *UserHandler) GetByID(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	err := ctx.ParamsParser(&req)
	if err != nil {
		return err
	}

	result, err := u.service.GetByID(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully get user by id",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

func (u *UserHandler) Update(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := u.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	claims := ctx.Locals("claims").(*domain.Claims)

	result, err := u.service.Update(req, claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully update user by id",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

func (u *UserHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	claims := ctx.Locals("claims").(*domain.Claims)

	err := u.service.Delete(req, claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully delete user by id",
	})
}
