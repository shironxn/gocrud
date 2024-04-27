package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"
)

type AuthHandler struct {
	service   port.AuthService
	jwt       util.JWT
	validator *util.Validator
	cfg       *config.Config
}

func NewAuthHandler(service port.AuthService, jwt util.JWT, validator *util.Validator, cfg *config.Config) port.AuthHandler {
	return &AuthHandler{
		service:   service,
		jwt:       jwt,
		validator: validator,
		cfg:       cfg,
	}
}

// @Summary Register a new user
// @Description Register a new user with the specified name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.AuthRegisterRequest true "User registration request object"
// @Success 201 {object} domain.UserResponse "Successfully registered a new user"
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req domain.AuthRegisterRequest

	cookie := ctx.Cookies("refresh-token")
	claims, _ := h.jwt.ValidateToken(cookie, h.cfg.JWT.Refresh)
	if claims != nil {
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
// @Param user body domain.AuthLoginRequest true "User login request object"
// @Success 200 {object} domain.UserResponse "Successfully logged in"
// @Router /auth/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req domain.AuthLoginRequest

	cookie := ctx.Cookies("refresh-token")
	claims, _ := h.jwt.ValidateToken(cookie, h.cfg.JWT.Refresh)
	if claims != nil {
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
		SameSite: func(dev string) string {
			if dev == "true" {
				return fiber.CookieSameSiteLaxMode
			}
			return fiber.CookieSameSiteNoneMode
		}(h.cfg.Server.Dev),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(10 * time.Minute),
		SameSite: func(dev string) string {
			if dev == "true" {
				return fiber.CookieSameSiteLaxMode
			}
			return fiber.CookieSameSiteNoneMode
		}(h.cfg.Server.Dev),
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "user successfully logged in",
		Data: domain.UserResponse{
			ID:        result.ID,
			Name:      result.Name,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
			UserToken: tokens,
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

// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags auth
// @Produce json
// @Success 200 {object} domain.SuccessResponse "Successfully refreshed token"
// @Router /auth/refresh [post]
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
		SameSite: func(dev string) string {
			if dev == "true" {
				return fiber.CookieSameSiteLaxMode
			}
			return fiber.CookieSameSiteNoneMode
		}(h.cfg.Server.Dev),
	})

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully refresh token",
		Data: domain.UserToken{
			AccessToken: *result,
			Claims:      claims,
		},
	})
}
