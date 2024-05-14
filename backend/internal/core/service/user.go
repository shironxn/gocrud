package service

import (
	"github.com/shironxn/blanknotes/internal/core/domain"
	"github.com/shironxn/blanknotes/internal/core/port"
	"github.com/shironxn/blanknotes/internal/util"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repository port.UserRepository
	bcrypt     util.Bcrypt
}

func NewUserService(repository port.UserRepository, bcrypt util.Bcrypt) port.UserService {
	return &UserService{
		repository: repository,
		bcrypt:     bcrypt,
	}
}

func (h *UserService) GetAll(req domain.UserQuery, metdata *domain.Metadata) ([]domain.User, error) {
	return h.repository.GetAll(req, metdata)
}

func (h *UserService) GetByID(id uint) (*domain.User, error) {
	return h.repository.GetByID(id)
}

func (h *UserService) Update(req domain.UserRequest, claims domain.Claims) (*domain.User, error) {
	user, err := h.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if user.ID != claims.UserID {
		return nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	if req.Password != "" {
		hashedPassword, err := h.bcrypt.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = string(hashedPassword)
	}

	return h.repository.Update(req, user)
}

func (h *UserService) Delete(req domain.UserRequest, claims domain.Claims) error {
	user, err := h.repository.GetByID(req.ID)
	if err != nil {
		return err
	}

	if user.ID != claims.UserID {
		return fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return h.repository.Delete(user)
}
