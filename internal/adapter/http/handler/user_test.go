package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gocrud/internal/config"
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"
	"gocrud/internal/mocks"
	"gocrud/internal/util"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserHandler_Register(t *testing.T) {
	type fields struct {
		service   port.UserService
		validator util.Validator
	}

	type args struct {
		req domain.UserRegisterRequest
	}

	expected := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserService.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(expected, nil).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     "shiron",
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully create user",
				Data: domain.UserResponse{
					ID:   expected.ID,
					Name: expected.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     "shiron",
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			code: fiber.StatusInternalServerError,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Post("/api/v1/auth/register", u.Register)

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
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	type fields struct {
		service   port.UserService
		validator util.Validator
		jwt       util.JWT
	}

	type args struct {
		req domain.UserLoginRequest
	}

	expected := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	}

	mockUserService := mocks.NewUserService(t)
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
				service: func() port.UserService {
					mockUserService.EXPECT().Login(mock.AnythingOfType("domain.UserLoginRequest")).Return(expected, nil).Once()
					return mockUserService
				}(),
				validator: validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully login user",
				Data: domain.UserResponse{
					ID:   expected.ID,
					Name: expected.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Login(mock.AnythingOfType("domain.UserLoginRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserService
				}(),
				validator: validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    "shiron@example.com",
					Password: "password123",
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
				service: func() port.UserService {
					mockUserService.EXPECT().Login(mock.AnythingOfType("domain.UserLoginRequest")).Return(nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")).Once()
					return mockUserService
				}(),
				validator: validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    "shiron@example.com",
					Password: "wrong",
				},
			},
			code: fiber.StatusUnauthorized,
			want: domain.ErrorResponse{
				Message: "invalid password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
				jwt:       tt.fields.jwt,
			}

			app := config.NewFiber()
			app.Post("/api/v1/auth/login", u.Login)

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
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestUserHandler_Logout(t *testing.T) {
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
				Value: "dummy-cookie",
			},
			want: domain.SuccessResponse{
				Message: "successfully logout user",
			},
			wantErr: false,
		},
		{
			name:   "failure",
			code:   fiber.StatusBadRequest,
			cookie: nil,
			want: domain.ErrorResponse{
				Message: "user already logged out",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{}

			app := config.NewFiber()
			app.Post("/api/v1/auth/register", u.Logout)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/register", nil)
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

func TestUserHandler_GetCurrent(t *testing.T) {
	type fields struct {
		service port.UserService
	}

	expected := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	}

	mockUserService := mocks.NewUserService(t)

	tests := []struct {
		name    string
		fields  fields
		code    int
		claims  *domain.Claims
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(expected, nil).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusOK,
			claims: &domain.Claims{
				UserID: expected.ID,
			},
			want: domain.SuccessResponse{
				Message: "successfully get current user data",
				Data: domain.UserResponse{
					ID:   expected.ID,
					Name: expected.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusInternalServerError,
			claims: &domain.Claims{
				UserID: expected.ID,
			},
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", tt.claims)
				return ctx.Next()
			})
			app.Get("/api/v1/user/current", u.GetCurrent)

			req := httptest.NewRequest(fiber.MethodGet, "/api/v1/user/current", nil)
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
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestUserHandler_GetAll(t *testing.T) {
	type fields struct {
		service port.UserService
	}

	expected := []domain.User{
		{
			Model:    gorm.Model{ID: 1},
			Name:     "shiron",
			Email:    "shiron@example.com",
			Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
		},
		{
			Model:    gorm.Model{ID: 2},
			Name:     "shironz",
			Email:    "shironz@example.com",
			Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
		},
		{
			Model:    gorm.Model{ID: 3},
			Name:     "shironzz",
			Email:    "shironzz@example.com",
			Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
		},
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
					mockUserService.EXPECT().GetAll().Return(expected, nil).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusOK,

			want: domain.SuccessResponse{
				Message: "successfully get all user data",
				Data: func() []domain.UserResponse {
					var data []domain.UserResponse
					for _, user := range expected {
						data = append(data, domain.UserResponse{
							ID:   user.ID,
							Name: user.Name,
						})
					}
					return data
				}(),
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().GetAll().Return(nil, errors.New("failed")).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusInternalServerError,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Get("/api/v1/user", u.GetAll)

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
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				for i, data := range got.Data.([]interface{}) {
					assert.Equal(t, tt.want.(domain.SuccessResponse).Data.([]domain.UserResponse)[i].ID, uint(data.(map[string]interface{})["id"].(float64)))
					assert.Equal(t, tt.want.(domain.SuccessResponse).Data.([]domain.UserResponse)[i].Name, data.(map[string]interface{})["name"].(string))
				}
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

	expected := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserService.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(expected, nil).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully get user by id",
				Data: domain.UserResponse{
					ID:   expected.ID,
					Name: expected.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusInternalServerError,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Get("/api/v1/user/:id", u.GetByID)

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
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	type fields struct {
		service   port.UserService
		validator util.Validator
	}

	type args struct {
		req domain.UserRequest
	}

	expected := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	}

	mockUserService := mocks.NewUserService(t)
	validator, _ := util.NewValidator()

	tests := []struct {
		name    string
		fields  fields
		args    args
		code    int
		claims  domain.Claims
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(expected, nil).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			args: args{
				domain.UserRequest{
					ID:       expected.ID,
					Name:     expected.Name,
					Email:    expected.Email,
					Password: expected.Password,
				},
			},
			code: fiber.StatusOK,
			claims: domain.Claims{
				UserID: expected.ID,
			},
			want: domain.SuccessResponse{
				Message: "successfully update user by id",
				Data: domain.UserResponse{
					ID:   expected.ID,
					Name: expected.Name,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(nil, errors.New("failed")).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			args: args{
				domain.UserRequest{
					ID:       expected.ID,
					Name:     expected.Name,
					Email:    expected.Email,
					Password: expected.Password,
				},
			},
			code: fiber.StatusInternalServerError,
			claims: domain.Claims{
				UserID: expected.ID,
			},
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
		{
			name: "permission",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), domain.Claims{UserID: expected.ID}).Return(nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockUserService
				}(),
				validator: validator,
			},
			args: args{
				domain.UserRequest{
					ID: func() uint {
						return uint(rand.Uint32())
					}(),
					Name:     expected.Name,
					Email:    expected.Email,
					Password: expected.Password,
				},
			},
			code: fiber.StatusForbidden,
			claims: domain.Claims{
				UserID: expected.ID,
			},
			want: domain.ErrorResponse{
				Message: "user does not have permission to perform this action",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", tt.claims)
				return ctx.Next()
			})
			app.Put("/api/v1/user/:id", u.Update)

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
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.UserResponse).Name, got.Data.(map[string]interface{})["name"].(string))
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	type fields struct {
		service port.UserService
	}

	type args struct {
		req domain.UserRequest
	}

	expected := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
	}

	mockUserService := mocks.NewUserService(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		code    int
		claims  domain.Claims
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
			args: args{
				domain.UserRequest{
					ID: expected.ID,
				},
			},
			code: fiber.StatusOK,
			claims: domain.Claims{
				UserID: expected.ID,
			},
			want: domain.SuccessResponse{
				Message: "successfully delete user by id",
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Delete(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(errors.New("failed")).Once()
					return mockUserService
				}(),
			},
			args: args{
				domain.UserRequest{
					ID: expected.ID,
				},
			},
			code: fiber.StatusInternalServerError,
			claims: domain.Claims{
				UserID: expected.ID,
			},
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
		{
			name: "permission",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Delete(mock.AnythingOfType("domain.UserRequest"), domain.Claims{UserID: expected.ID}).Return(fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockUserService
				}(),
			},
			args: args{
				domain.UserRequest{
					ID: func() uint {
						return uint(rand.Uint32())
					}(),
				},
			},
			code: fiber.StatusForbidden,
			claims: domain.Claims{
				UserID: expected.ID,
			},
			want: domain.ErrorResponse{
				Message: "user does not have permission to perform this action",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", tt.claims)
				return ctx.Next()
			})
			app.Delete("/api/v1/user/:id", u.Delete)

			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/v1/user/%v", tt.args.req.ID), nil)
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
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
			}
		})
	}
}
