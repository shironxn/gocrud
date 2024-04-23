package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"
)

type AuthService struct {
	repository port.AuthRepository
	bcrypt     util.Bcrypt
	jwt        util.JWT
	cfg        *config.Config
}

func NewAuthService(repository port.AuthRepository, bcrypt util.Bcrypt, jwt util.JWT, cfg *config.Config) port.AuthService {
	return &AuthService{
		repository: repository,
		bcrypt:     bcrypt,
		jwt:        jwt,
		cfg:        cfg,
	}
}

func (s *AuthService) Register(req domain.AuthRegisterRequest) (*domain.User, error) {
	hashedPassword, err := s.bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = string(hashedPassword)

	return s.repository.Register(req)
}

func (s *AuthService) Login(req domain.AuthLoginRequest) (*domain.User, *domain.UserToken, error) {
	user, err := s.repository.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}

	if err := s.bcrypt.ComparePassword(req.Password, []byte(user.Password)); err != nil {
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repository.StoreRefreshToken(user.ID, refreshToken); err != nil {
		return nil, nil, err
	}

	return user, &domain.UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Claims:       nil,
	}, nil
}

func (s *AuthService) Logout(userID uint) error {
	refresh, err := s.repository.GetRefreshToken(userID)
	if err != nil {
		return err
	}

	return s.repository.DeleteRefreshToken(*refresh)
}

func (s *AuthService) Refresh(token string) (*string, *domain.Claims, error) {
	claims, err := s.jwt.ValidateToken(token, s.cfg.JWT.Refresh)
	if err != nil {
		return nil, nil, err
	}

	refresh, err := s.repository.GetRefreshToken(claims.UserID)
	if err != nil {
		return nil, nil, err
	}

	if refresh.Token != token {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	accessToken, err := s.jwt.GenerateAccessToken(claims.UserID)
	if err != nil {
		return nil, nil, err
	}

	return &accessToken, claims, nil
}
