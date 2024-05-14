package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shironxn/blanknotes/internal/config"
	"github.com/shironxn/blanknotes/internal/core/domain"
	"github.com/shironxn/blanknotes/internal/core/port"
	"github.com/shironxn/blanknotes/internal/mocks"
	"github.com/shironxn/blanknotes/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var authEntity = domain.User{
	Model:     gorm.Model{ID: 1},
	Name:      "shiron",
	Email:     "shiron@example.com",
	Bio:       "hello world",
	AvatarURL: "https://i.pinimg.com/originals/be/38/3b/be383bedd646e4dd8a8e7c0cc304f9e9.jpg",
	Password:  "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	Notes:     []domain.Note{},
}

var userToken = domain.UserToken{
	AccessToken:  "anjay",
	RefreshToken: "anjay",
}

func TestAuthHandler_Register(t *testing.T) {
	type fields struct {
		service   port.AuthService
		validator *util.Validator
	}

	type args struct {
		req domain.AuthRegisterRequest
	}

	mockAuthService := mocks.NewAuthService(t)
	validator, _ := util.NewValidator()

	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Register(mock.AnythingOfType("domain.AuthRegisterRequest")).Return(&authEntity, nil).Once()
					return mockAuthService
				}(),
				validator: validator,
			},
			args: args{
				req: domain.AuthRegisterRequest{
					Name:     authEntity.Name,
					Email:    authEntity.Email,
					Password: authEntity.Password,
				},
			},
			code: fiber.StatusCreated,
		},
		{
			name: "validation error",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Register(mock.AnythingOfType("domain.AuthRegisterRequest")).Maybe()
					return mockAuthService
				}(),
				validator: validator,
			},
			code: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Post("/api/v1/auth/register", h.Register)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type fields struct {
		service   port.AuthService
		cfg       *config.Config
		jwt       util.JWT
		validator *util.Validator
	}

	type args struct {
		req domain.AuthLoginRequest
	}

	mockAuthService := mocks.NewAuthService(t)
	jwt := util.NewJWT(&config.Config{})
	validator, _ := util.NewValidator()

	tests := []struct {
		name    string
		fields  fields
		args    args
		code    int
		wantErr interface{}
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Login(mock.AnythingOfType("domain.AuthLoginRequest")).Return(&authEntity, &userToken, nil).Once()
					return mockAuthService
				}(),
				cfg:       &config.Config{},
				jwt:       jwt,
				validator: validator,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: authEntity.Password,
				},
			},
			code: fiber.StatusOK,
		},
		{
			name: "wrong password",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Login(mock.AnythingOfType("domain.AuthLoginRequest")).Return(nil, nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")).Once()
					return mockAuthService
				}(),
				jwt:       jwt,
				validator: validator,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: authEntity.Password + "lol",
				},
			},
			code: fiber.StatusUnauthorized,
			wantErr: domain.ErrorResponse{
				Code:  401,
				Error: "invalid password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{
				cfg:       tt.fields.cfg,
				service:   tt.fields.service,
				jwt:       tt.fields.jwt,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Post("/api/v1/auth/login", h.Login)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)

			if tt.wantErr != nil {
				var got domain.ErrorResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantErr, got)
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	type fields struct {
		service port.AuthService
		jwt     util.JWT
	}

	type args struct {
		claims domain.Claims
	}

	mockAuthService := mocks.NewAuthService(t)
	jwt := util.NewJWT(&config.Config{})

	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
		cookie []*fiber.Cookie
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Logout(mock.AnythingOfType("uint")).Return(nil)
					return mockAuthService
				}(),
				jwt: jwt,
			},
			args: args{
				domain.Claims{
					UserID: authEntity.ID,
				},
			},
			code: fiber.StatusOK,
			cookie: []*fiber.Cookie{
				{
					Name:  "refresh-token",
					Value: "dummy-token",
				},
				{
					Name:  "access-token",
					Value: "dummy-token",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{
				service: tt.fields.service,
				jwt:     tt.fields.jwt,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Post("/api/v1/auth/logout", h.Logout)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.cookie != nil {
				for _, cookie := range tt.cookie {
					req.AddCookie(&http.Cookie{
						Name:  cookie.Name,
						Value: cookie.Value,
					})
				}
			}

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
