package util

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"gocrud/internal/config"
	"gocrud/internal/core/domain"
)

type JWT struct {
	cfg *config.Config
}

func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		cfg: cfg,
	}
}

func (j *JWT) GenerateToken(ctx *fiber.Ctx, user *domain.User) (*domain.UserDetails, error) {
	jwtExpire := time.Now().Add(1 * time.Hour)

	claims := domain.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwtExpire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	cookieExpire := time.Now().Add(1 * time.Hour)
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HTTPOnly: true,
		Expires:  cookieExpire,
	})

	details := domain.UserDetails{
		Token:     tokenString,
		ExpiredAt: cookieExpire.String(),
	}

	return &details, nil
}

func (j *JWT) ValidateToken(token string) (jwt.Claims, error) {
	tokenString, err := jwt.ParseWithClaims(token, &domain.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !tokenString.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := tokenString.Claims.(*domain.Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
