package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/mocks"
	"github.com/shironxn/gocrud/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var authEntity = &domain.User{
	Model:     gorm.Model{ID: 1},
	Name:      "shiron",
	Email:     "shiron@example.com",
	Bio:       "hello world",
	AvatarURL: "https://i.pinimg.com/originals/be/38/3b/be383bedd646e4dd8a8e7c0cc304f9e9.jpg",
	Password:  "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	Notes:     []domain.Note{},
}

var userToken = &domain.UserToken{
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
		name    string
		fields  fields
		args    args
		code    int
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Register(mock.AnythingOfType("domain.AuthRegisterRequest")).Return(authEntity, nil).Once()
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
			want: domain.SuccessResponse{
				Message: "user successfully registered",
				Data: domain.UserResponse{
					ID:   authEntity.ID,
					Name: authEntity.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Register(mock.AnythingOfType("domain.AuthRegisterRequest")).Return(nil, errors.New("failed")).Once()
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
			code: fiber.StatusInternalServerError,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
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
			want: domain.ErrorResponse{
				Message: "validation error",
			},
			wantErr: true,
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

			if tt.wantErr {
				var got domain.ErrorResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Name, authEntity.Name)
				assert.Equal(t, tt.args.req.Email, authEntity.Email)
				assert.Equal(t, tt.args.req.Password, authEntity.Password)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type fields struct {
		service   port.AuthService
		validator *util.Validator
		jwt       util.JWT
	}

	type args struct {
		req domain.AuthLoginRequest
	}

	mockAuthService := mocks.NewAuthService(t)
	validator, _ := util.NewValidator()
	jwt := util.NewJWT(&config.Config{})

	tests := []struct {
		name    string
		fields  fields
		args    args
		code    int
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Login(mock.AnythingOfType("domain.AuthLoginRequest")).Return(authEntity, userToken, nil).Once()
					return mockAuthService
				}(),
				validator: validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: authEntity.Password,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "user successfully logged in",
				Data: domain.UserResponse{
					ID:   authEntity.ID,
					Name: authEntity.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Login(mock.AnythingOfType("domain.AuthLoginRequest")).Return(nil, nil, errors.New("failed")).Once()
					return mockAuthService
				}(),
				validator: validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: authEntity.Password,
				},
			},
			code: fiber.StatusInternalServerError,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
		{
			name: "wrong password",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Login(mock.AnythingOfType("domain.AuthLoginRequest")).Return(nil, nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")).Once()
					return mockAuthService
				}(),
				validator: validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: authEntity.Password + "lol",
				},
			},
			code: fiber.StatusUnauthorized,
			want: domain.ErrorResponse{
				Message: "invalid password",
			},
			wantErr: true,
		},
		{
			name: "validation error",
			fields: fields{
				service: func() port.AuthService {
					mockAuthService.EXPECT().Login(mock.AnythingOfType("domain.AuthLoginRequest")).Maybe()
					return mockAuthService
				}(),
				validator: validator,
			},
			code: fiber.StatusBadRequest,
			want: domain.ErrorResponse{
				Message: "validation error",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
				jwt:       tt.fields.jwt,
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

			if tt.wantErr {
				var got domain.ErrorResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Email, authEntity.Email)
				assert.Equal(t, tt.args.req.Password, authEntity.Password)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		cookie  *http.Cookie
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			code: fiber.StatusOK,
			cookie: &http.Cookie{
				Name:  "token",
				Value: "dummy-token",
			},
			want: domain.SuccessResponse{
				Message: "user successfully logged out",
			},
			wantErr: false,
		},
		{
			name:   "failure",
			code:   fiber.StatusBadRequest,
			cookie: nil,
			want: domain.ErrorResponse{
				Message: "user is already logged out",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{}

			app := config.NewFiber()
			app.Post("/api/v1/auth/logout", h.Logout)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)

			if tt.wantErr {
				var got domain.ErrorResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)

				var tokenCookie *http.Cookie
				for _, cookie := range res.Cookies() {
					if cookie.Name == tt.cookie.Name {
						tokenCookie = cookie
						break
					}
				}

				assert.NotNil(t, tokenCookie)
				assert.Empty(t, tokenCookie.Value)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
			}
		})
	}
}
