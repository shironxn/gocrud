package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestUserHandler_Register(t *testing.T) {
	type fields struct {
		service   port.UserService
		validator util.Validator
	}

	type args struct {
		req domain.UserRegisterRequest
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
					mockUserService.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     userEntity.Name,
					Email:    userEntity.Email,
					Password: userEntity.Password,
				},
			},
			code: fiber.StatusCreated,
			want: domain.SuccessResponse{
				Message: "user successfully registered",
				Data: domain.UserResponse{
					ID:   userEntity.ID,
					Name: userEntity.Name,
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
				validator: *validator,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     userEntity.Name,
					Email:    userEntity.Email,
					Password: userEntity.Password,
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
				service: func() port.UserService {
					mockUserService.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Maybe()
					return mockUserService
				}(),
				validator: *validator,
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
			h := &UserHandler{
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
				assert.Equal(t, tt.args.req.Name, userEntity.Name)
				assert.Equal(t, tt.args.req.Email, userEntity.Email)
				assert.Equal(t, tt.args.req.Password, userEntity.Password)
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
					mockUserService.EXPECT().Login(mock.AnythingOfType("domain.UserLoginRequest")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
				validator: *validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    userEntity.Email,
					Password: userEntity.Password,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "user successfully logged in",
				Data: domain.UserResponse{
					ID:   userEntity.ID,
					Name: userEntity.Name,
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
				validator: *validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    userEntity.Email,
					Password: userEntity.Password,
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
				validator: *validator,
				jwt:       jwt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    userEntity.Email,
					Password: userEntity.Password + "lol",
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
				service: func() port.UserService {
					mockUserService.EXPECT().Login(mock.AnythingOfType("domain.UserLoginRequest")).Maybe()
					return mockUserService
				}(),
				validator: *validator,
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
			h := &UserHandler{
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
				assert.Equal(t, tt.args.req.Email, userEntity.Email)
				assert.Equal(t, tt.args.req.Password, userEntity.Password)
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
			h := &UserHandler{}

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
					mockUserService.EXPECT().GetAll(mock.AnythingOfType("domain.UserQuery"), mock.AnythingOfType("domain.Metadata")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully retrieved all user data",
				Data: func() []domain.UserResponse {
					var data []domain.UserResponse
					for _, user := range userEntity {
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
					mockUserService.EXPECT().GetAll(mock.AnythingOfType("domain.UserQuery"), mock.AnythingOfType("domain.Metadata")).Return(nil, errors.New("failed")).Once()
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
					mockUserService.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully retrieved user by id",
				Data: domain.UserResponse{
					ID:   userEntity.ID,
					Name: userEntity.Name,
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
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, userEntity.ID)
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
		req    domain.UserRequest
		claims domain.Claims
	}

	userEntity := &domain.User{
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
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(userEntity, nil).Once()
					return mockUserService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.UserRequest{
					ID:       userEntity.ID,
					Name:     userEntity.Name,
					Email:    userEntity.Email,
					Password: userEntity.Password,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully updated user by id",
				Data: domain.UserResponse{
					ID:   userEntity.ID,
					Name: userEntity.Name,
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
				validator: *validator,
			},
			args: args{
				req: domain.UserRequest{
					ID:       userEntity.ID,
					Name:     userEntity.Name,
					Email:    userEntity.Email,
					Password: userEntity.Password,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			code: fiber.StatusInternalServerError,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockUserService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.UserRequest{
					ID:       userEntity.ID,
					Name:     userEntity.Name,
					Email:    userEntity.Email,
					Password: userEntity.Password,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			code: fiber.StatusForbidden,
			want: domain.ErrorResponse{
				Message: "user does not have permission to perform this action",
			},
			wantErr: true,
		},
		{
			name: "validation error",
			fields: fields{
				service:   mockUserService,
				validator: *validator,
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
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, userEntity.ID)
				assert.Equal(t, tt.args.req.Name, userEntity.Name)
				assert.Equal(t, tt.args.req.Email, userEntity.Email)
				assert.Equal(t, tt.args.req.Password, userEntity.Password)
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
		req    domain.UserRequest
		claims domain.Claims
	}

	userEntity := &domain.User{
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
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			code: fiber.StatusOK,
			cookie: &http.Cookie{
				Name:  "token",
				Value: "dummy-token",
			},
			want: domain.SuccessResponse{
				Message: "successfully deleted user by id",
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
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			code:   http.StatusInternalServerError,
			cookie: nil,
			want: domain.ErrorResponse{
				Message: "failed",
			},
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				service: func() port.UserService {
					mockUserService.EXPECT().Delete(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("domain.Claims")).Return(fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockUserService
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			code:   http.StatusForbidden,
			cookie: nil,
			want: domain.ErrorResponse{
				Message: "user does not have permission to perform this action",
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
				assert.Equal(t, tt.args.req.ID, userEntity.ID)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
			}
		})
	}
}
