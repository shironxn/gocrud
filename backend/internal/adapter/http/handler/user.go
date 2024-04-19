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
	validator *util.Validator
	jwt       util.JWT
}

func NewUserHandler(service port.UserService, validator *util.Validator, jwt util.JWT) port.UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator,
		jwt:       jwt,
	}
}

// @Summary Get all users
// @Description Retrieve data of all registered users
// @Tags user
// @Produce json
// @Param id query int false "Filter users by ID"
// @Param name query string false "Filter users by name"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} []domain.UserPaginationResponse "Successfully retrieved all user data"
// @Router /users [get]
func (h *UserHandler) GetAll(ctx *fiber.Ctx) error {
	var req domain.UserQuery
	var metadata domain.Metadata
	var data []domain.UserResponse

	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

	if err := ctx.QueryParser(&metadata); err != nil {
		return err
	}

	result, err := h.service.GetAll(req, &metadata)
	if err != nil {
		return err
	}

	for _, user := range result {
		data = append(data,
			domain.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Bio:       user.Bio,
				AvatarURL: user.AvatarURL,
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
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} domain.UserResponse "Successfully retrieved user by ID"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.GetByID(req.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully retrieved user by id",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			Bio:       result.Bio,
			AvatarURL: result.AvatarURL,
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
// @Router /users/{id} [put]
func (h *UserHandler) Update(ctx *fiber.Ctx) error {
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

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := h.service.Update(req, *claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully updated user by id",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			Bio:       result.Bio,
			AvatarURL: result.AvatarURL,
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
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.UserRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}

	err := h.service.Delete(req, *claims)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully deleted user by id",
	})
}
