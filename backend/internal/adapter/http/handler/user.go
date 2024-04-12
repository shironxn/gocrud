package handler

import (
	"reflect"
	"time"

	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

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
		}},
	)
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

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "user successfully logged out",
	})
}

// @Summary Get all users
// @Description Retrieve data of all registered users
// @Tags user
// @Produce json
// @Param id query int false "Filter users by ID"
// @Param name query string false "Filter users by name"
// @Success 200 {object} []domain.UserPaginationResponse "Successfully retrieved all user data"
// @Router /users [get]
func (u *UserHandler) GetAll(ctx *fiber.Ctx) error {
	var req domain.UserQuery
	var metadata domain.Metadata
	var data []domain.UserResponse

	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

	if err := ctx.QueryParser(&metadata); err != nil {
		return err
	}

	result, err := u.service.GetAll(req, &metadata)
	if err != nil {
		return err
	}

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
		Data: domain.UserPaginationResponse{
			Users:    data,
			Metadata: metadata,
		},
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

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
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

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if reflect.DeepEqual(req, domain.UserRequest{}) {
		return fiber.NewError(fiber.StatusBadRequest, "at least one field must be filled")
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}

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
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}

	err := u.service.Delete(req, *claims)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully deleted user by id",
	})
}
