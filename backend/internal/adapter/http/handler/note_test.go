package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var noteEntity = &domain.Note{
	Model: gorm.Model{
		ID: 1,
	},
	Title:       "golang",
	Description: "lets go",
	CoverURL:    "https://i.pinimg.com/originals/56/c3/ee/56c3ee9cae0c8152bd341b969cd2fc1d.png",
	Content:     "is the best",
	Visibility:  "public",
	UserID:      1,
}

// var metadataEntity = domain.Metadata{
// 	Sort:         "id",
// 	Order:        "desc",
// 	TotalRecords: 100,
// 	TotalPage:    10,
// 	Limit:        10,
// 	Page:         1,
// }

func TestNoteHandler_Create(t *testing.T) {
	type fields struct {
		service   port.NoteService
		validator *util.Validator
	}

	type args struct {
		req    domain.NoteRequest
		claims domain.Claims
	}

	mockNoteService := mocks.NewNoteService(t)
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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
				validator: validator,
			},
			args: args{
				req: domain.NoteRequest{
					Title:       noteEntity.Title,
					Description: noteEntity.Description,
					CoverURL:    noteEntity.CoverURL,
					Content:     noteEntity.Content,
					Visibility:  string(noteEntity.Visibility),
				},
			},
			code: fiber.StatusCreated,
		},
		{
			name: "validation error",
			fields: fields{
				service:   mockNoteService,
				validator: validator,
			},
			code: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Post("/api/v1/notes", h.Create)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/notes", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func TestNoteHandler_GetAll(t *testing.T) {
	type fields struct {
		service port.NoteService
		jwt     util.JWT
	}

	mockNoteService := mocks.NewNoteService(t)
	jwt := util.NewJWT(&config.Config{})

	tests := []struct {
		name   string
		fields fields
		code   int
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().GetAll(mock.AnythingOfType("domain.NoteQuery"), mock.AnythingOfType("*domain.Metadata")).Return([]domain.Note{
						*noteEntity,
					}, nil).Once()
					return mockNoteService
				}(),
				jwt: jwt,
			},
			code: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteHandler{
				service: tt.fields.service,
				jwt:     tt.fields.jwt,
			}

			app := config.NewFiber()
			app.Get("/api/v1/notes", h.GetAll)

			req := httptest.NewRequest(fiber.MethodGet, "/api/v1/notes", nil)
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func TestNoteHandler_GetByID(t *testing.T) {
	type fields struct {
		service port.NoteService
		jwt     util.JWT
	}

	type args struct {
		req    uint
		claims *domain.Claims
	}

	mockNoteService := mocks.NewNoteService(t)
	jwt := util.NewJWT(&config.Config{})

	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().GetByID(mock.AnythingOfType("uint"), mock.AnythingOfType("*domain.Claims")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
				jwt: jwt,
			},
			code: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteHandler{
				service: tt.fields.service,
				jwt:     tt.fields.jwt,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", tt.args.claims)
				return ctx.Next()
			})
			app.Get("/api/v1/note/:id", h.GetByID)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/v1/note/%v", tt.args.req), bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func TestNoteHandler_Update(t *testing.T) {
	type fields struct {
		service   port.NoteService
		jwt       util.JWT
		validator *util.Validator
	}

	type args struct {
		req    domain.NoteUpdateRequest
		claims domain.Claims
	}

	mockNoteService := mocks.NewNoteService(t)
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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Update(mock.AnythingOfType("domain.NoteUpdateRequest"), mock.AnythingOfType("domain.Claims")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
				jwt:       jwt,
				validator: validator,
			},
			code: fiber.StatusOK,
		},
		{
			name: "permission denied",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Update(mock.AnythingOfType("domain.NoteUpdateRequest"), mock.AnythingOfType("domain.Claims")).Return(nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockNoteService
				}(),
				validator: validator,
			},
			code: fiber.StatusForbidden,
			wantErr: domain.ErrorResponse{
				Code:  403,
				Error: "user does not have permission to perform this action",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteHandler{
				service:   tt.fields.service,
				jwt:       tt.fields.jwt,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Put("/api/v1/note/:id", h.Update)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/api/v1/note/%v", tt.args.req.ID), bytes.NewBuffer(requestBody))
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

func TestNoteHandler_Delete(t *testing.T) {
	type fields struct {
		service port.NoteService
	}

	type args struct {
		req    uint
		claims domain.Claims
	}

	mockNoteService := mocks.NewNoteService(t)

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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Delete(mock.AnythingOfType("uint"), mock.AnythingOfType("domain.Claims")).Return(nil).Once()
					return mockNoteService
				}(),
			},
			code: fiber.StatusOK,
		},
		{
			name: "permission denied",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Delete(mock.AnythingOfType("uint"), mock.AnythingOfType("domain.Claims")).Return(fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockNoteService
				}(),
			},
			code: fiber.StatusForbidden,
			wantErr: domain.ErrorResponse{
				Code:  403,
				Error: "user does not have permission to perform this action",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Delete("/api/v1/note/:id", h.Delete)

			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/v1/note/%v", tt.args.req), nil)
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
