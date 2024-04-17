package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"
)

type AuthHandler struct {
	service   port.AuthService
	jwt       util.JWT
	validator util.Validator
}

func NewAuthHandler(service port.AuthService, jwt util.JWT, validator util.Validator) port.AuthHandler {
	return &AuthHandler{
		service:   service,
		jwt:       jwt,
		validator: validator,
	}
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req domain.UserRegisterRequest

	if cookie := ctx.Cookies("token"); cookie != "" {
		return fiber.NewError(fiber.StatusBadRequest, "user is already registered")
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := h.service.Register(req)
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
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req domain.UserLoginRequest

	if cookie := ctx.Cookies("token"); cookie != "" {
		return fiber.NewError(fiber.StatusBadRequest, "user is already logged in")
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, tokens, err := h.service.Login(req)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		// SameSite: fiber.CookieSameSiteNoneMode,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(10 * time.Minute),
		// SameSite: fiber.CookieSameSiteNoneMode,
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "user successfully logged in",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
			UserToken: *tokens,
		}},
	)
}

// @Summary User logout
// @Description Log out the currently logged-in user
// @Tags auth
// @Produce json
// @Success 200 {object} domain.SuccessResponse "Successfully logged out"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	var req domain.User

	// if cookie := ctx.Cookies("token"); cookie == "" {
	// 	return fiber.NewError(fiber.StatusBadRequest, "user is already logged out")
	// }

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
	req.ID = uint(claims.UserID)

	if err := h.service.Logout(req.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "user successfully logged out",
	})
}

func (h *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh-token")
	result, claims, err := h.service.Refresh(cookie)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Value:    *result,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(10 * time.Minute),
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully refresh token",
		Data: domain.UserToken{
			AccessToken: *result,
			Claims:      claims,
		},
	})
}
