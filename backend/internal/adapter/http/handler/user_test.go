package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/mocks"
	"github.com/shironxn/gocrud/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var userEntity = &domain.User{
	Model:     gorm.Model{ID: 1},
	Name:      "shiron",
	Email:     "shiron@example.com",
	Bio:       "hello world",
	AvatarURL: "https://i.pinimg.com/originals/be/38/3b/be383bedd646e4dd8a8e7c0cc304f9e9.jpg",
	Password:  "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	Notes:     []domain.Note{},
}

func TestUserHandler_GetAll(t *testing.T) {
	type fields struct {
		service port.UserService
	}

	userEntity := []domain.User{
		*userEntity,
	}

	mockUserService := mocks.NewUserService(t)

	tests := []struct {
		name    string
		fields  fields
		code    int
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().GetAll(mock.AnythingOfType("domain.UserQuery"), mock.AnythingOfType("*domain.Metadata")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
			},
			code:    fiber.StatusOK,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Get("/api/v1/user", h.GetAll)

			req := httptest.NewRequest(fiber.MethodGet, "/api/v1/user", nil)
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
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	type fields struct {
		service port.UserService
	}

	type args struct {
		req domain.UserRequest
	}

	mockUserService := mocks.NewUserService(t)

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
				service: func() port.UserService {
					mockUserService.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
			},
			code:    fiber.StatusOK,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Get("/api/v1/user/:id", h.GetByID)

			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/v1/user/%v", tt.args.req.ID), nil)
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
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	type fields struct {
		service   port.UserService
		validator *util.Validator
	}

	type args struct {
		req    domain.UserRequest
		claims domain.Claims
	}

	mockUserService := mocks.NewUserService(t)
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
				service: func() port.UserService {
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			code:    fiber.StatusOK,
			wantErr: false,
		},
		{
			name: "permission denied",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			code: fiber.StatusForbidden,
			want: domain.ErrorResponse{
				Code:  403,
				Error: "user does not have permission to perform this action",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Put("/api/v1/user/:id", h.Update)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/api/v1/user/%v", tt.args.req.ID), bytes.NewBuffer(requestBody))
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
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	type fields struct {
		service port.UserService
	}

	type args struct {
		req    domain.UserRequest
		claims domain.Claims
	}

	mockUserService := mocks.NewUserService(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		code    int
		cookie  *http.Cookie
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Delete(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(nil).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusOK,
			cookie: &http.Cookie{
				Name:  "refresh-token",
				Value: "dummy-token",
			},
			wantErr: false,
		},
		{
			name: "permission denied",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Delete(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockUserService
				}(),
			},
			code:   http.StatusForbidden,
			cookie: nil,
			want: domain.ErrorResponse{
				Code:  403,
				Error: "user does not have permission to perform this action",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Delete("/api/v1/user/:id", h.Delete)

			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/v1/user/%v", tt.args.req.ID), nil)
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
				var tokenCookie *http.Cookie
				for _, cookie := range res.Cookies() {
					if cookie.Name == tt.cookie.Name {
						tokenCookie = cookie
						break
					}
				}

				assert.NoError(t, err)
				assert.NotNil(t, tokenCookie)
				assert.Empty(t, tokenCookie.Value)
			}
		})
	}
}
