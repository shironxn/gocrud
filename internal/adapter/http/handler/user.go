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
	validator util.Validator
	jwt       util.JWT
}

func NewUserHandler(service port.UserService, validator util.Validator, jwt util.JWT) port.UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator,
		jwt:       jwt,
	}
}

// @Summary Register a new user
// @Description Register a new user with the specified name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.UserRegisterRequest true "User registration request object"
// @Success 201 {object} domain.UserResponse "Successfully registered a new user"
// @Failure 400 {object} domain.ErrorValidationResponse "Validation error"
// @Router /auth/register [post]
func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	var req domain.UserRegisterRequest

	if cookie := ctx.Cookies("token"); cookie != "" {
		return fiber.NewError(fiber.StatusBadRequest, "user is already registered")
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

	return ctx.Status(fiber.StatusCreated).JSON(domain.SuccessResponse{
		Message: "user successfully registered",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

// @Summary User login
// @Description Log in an existing user with the provided email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.UserLoginRequest true "User login request object"
// @Success 200 {object} domain.UserResponse "Successfully logged in"
// @Failure 400 {object} domain.ErrorValidationResponse "Validation error"
// @Router /auth/login [post]
func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	var req domain.UserLoginRequest

	if cookie := ctx.Cookies("token"); cookie != "" {
		return fiber.NewError(fiber.StatusBadRequest, "user is already logged in")
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
		Message: "user successfully logged in",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}})
}

// @Summary User logout
// @Description Log out the currently logged-in user
// @Tags auth
// @Produce json
// @Success 200 {object} domain.SuccessResponse "Successfully logged out"
// @Router /auth/logout [post]
func (u *UserHandler) Logout(ctx *fiber.Ctx) error {
	if cookie := ctx.Cookies("token"); cookie == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user is already logged out")
	}

	cookieExpire := time.Now().Add(-time.Hour * 24)
	ctx.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: cookieExpire,
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "user successfully logged out",
	})
}

// @Summary Get current user data
// @Description Retrieve data of the currently logged-in user
// @Tags user
// @Produce json
// @Success 200 {object} domain.UserResponse "Successfully retrieved current user data"
// @Router /user/current [get]
func (u *UserHandler) GetCurrent(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	claims := ctx.Locals("claims").(*domain.Claims)
	req.ID = claims.UserID

	result, err := u.service.GetByID(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully retrieved current user data",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

// @Summary Get all users
// @Description Retrieve data of all registered users
// @Tags user
// @Produce json
// @Success 200 {object} []domain.UserResponse "Successfully retrieved all user data"
// @Failure 400 {object} domain.ErrorValidationResponse "Validation error"
// @Router /user [get]
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
		Message: "successfully retrieved all user data",
		Data:    data,
	})
}

// @Summary Get a user by ID
// @Description Retrieve data of a user based on the provided ID
// @Tags user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.UserResponse "Successfully retrieved user by ID"
// @Router /user/{id} [get]
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
		Message: "successfully retrieved user by id",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

// @Summary Update user data by ID
// @Description Update data of an existing user based on the provided ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.UserRequest true "Updated user data object"
// @Success 200 {object} domain.UserResponse "Successfully updated user by ID"
// @Failure 400 {object} domain.ErrorValidationResponse "Validation error"
// @Router /user/{id} [put]
func (u *UserHandler) Update(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims := ctx.Locals("claims").(*domain.Claims)

	if err := u.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := u.service.Update(req, *claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully updated user by id",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

// @Summary Delete a user by ID
// @Description Delete an existing user based on the provided ID
// @Tags user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.SuccessResponse "Successfully deleted user by ID"
// @Router /user/{id} [delete]
func (u *UserHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	claims := ctx.Locals("claims").(*domain.Claims)

	err := u.service.Delete(req, *claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully deleted user by id",
	})
}
