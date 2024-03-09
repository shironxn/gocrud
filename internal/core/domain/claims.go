package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint `json:"user_id" cookie:"user_id"`
	jwt.RegisteredClaims
}
