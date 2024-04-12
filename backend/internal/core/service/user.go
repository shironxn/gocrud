package service

import (
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

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

func (u *UserService) Create(req domain.UserRegisterRequest) (*domain.User, error) {
	hashedPassword, err := u.bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = string(hashedPassword)

	return u.repository.Create(req)
}

func (u *UserService) Login(req domain.UserLoginRequest) (*domain.User, error) {
	data, err := u.repository.GetByEmail(domain.UserRequest{Email: req.Email})
	if err != nil {
		return nil, err
	}

	if err := u.bcrypt.ComparePassword(req.Password, []byte(data.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	return data, nil
}

func (u *UserService) GetAll(req domain.UserQuery, metdata *domain.Metadata) ([]domain.User, error) {
	return u.repository.GetAll(req, metdata)
}

func (u *UserService) GetByID(req domain.UserRequest) (*domain.User, error) {
	return u.repository.GetByID(req)
}

func (u *UserService) Update(req domain.UserRequest, claims domain.Claims) (*domain.User, error) {
	user, err := u.repository.GetByID(req)
	if err != nil {
		return nil, err
	}

	if user.ID != claims.UserID {
		return nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	if req.Password != "" {
		hashedPassword, err := u.bcrypt.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = string(hashedPassword)
	}

	return u.repository.Update(req, user)
}

func (u *UserService) Delete(req domain.UserRequest, claims domain.Claims) error {
	user, err := u.repository.GetByID(req)
	if err != nil {
		return err
	}

	if user.ID != claims.UserID {
		return fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return u.repository.Delete(user)
}
