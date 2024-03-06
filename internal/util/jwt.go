package util

import (
	"errors"
	"gocrud/internal/config"
	"gocrud/internal/core/domain"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
}

func (j *JWTManager) GenerateToken(w http.ResponseWriter, user domain.User) (*domain.UserDetails, error) {
	config := config.GetConfig()

	jwtExpire := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))

	claims := domain.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwtExpire,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))

	cookieExpire := time.Now().Add(1 * time.Hour)
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  cookieExpire,
	}

	http.SetCookie(w, cookie)

	details := domain.UserDetails{
		Token:     tokenString,
		ExpiredAt: cookieExpire.String(),
	}

	return &details, err
}

func (j *JWTManager) ValidateToken(token string) (jwt.Claims, error) {
	config := config.GetConfig()

	tokenString, err := jwt.ParseWithClaims(token, &domain.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	log.Info(tokenString)

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

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
