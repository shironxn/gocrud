package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
)

type JWT struct {
	cfg *config.Config
}

func NewJWT(cfg *config.Config) JWT {
	return JWT{
		cfg: cfg,
	}
}

func (j JWT) GenerateAccessToken(userID uint) (*string, error) {
	accessTokenClaims := domain.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(j.cfg.JWT.Access))
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (j JWT) GenerateRefreshToken(userID uint) (*string, error) {
	refreshTokenClaims := domain.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(j.cfg.JWT.Refresh))
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (j JWT) ValidateToken(token string, secret string) (*domain.Claims, error) {
	tokenString, err := jwt.ParseWithClaims(token, &domain.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	if !tokenString.Valid {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	claims, ok := tokenString.Claims.(*domain.Claims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid token claims")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, fiber.NewError(fiber.StatusBadRequest, "token has expired")
	}

	return claims, nil
}
