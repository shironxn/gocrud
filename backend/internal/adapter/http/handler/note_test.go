package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Title:      "golang",
	Content:    "is the best",
	Visibility: "public",
	UserID:     1,
}

var metadataEntity = domain.Metadata{
	Sort:         "id",
	Order:        "desc",
	TotalRecords: 100,
	TotalPage:    10,
	Limit:        10,
	Page:         1,
}

func TestNoteHandler_Create(t *testing.T) {
	type fields struct {
		service   port.NoteService
		validator util.Validator
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
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.NoteRequest{
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
				},
				claims: domain.Claims{
					UserID: noteEntity.ID,
				},
			},
			code: fiber.StatusCreated,
			want: domain.SuccessResponse{
				Message: "successfully created note",
				Data: domain.NoteResponse{
					ID:         noteEntity.ID,
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
					UserID:     noteEntity.UserID,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(nil, errors.New("failed")).Once()
					return mockNoteService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.NoteRequest{
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
				},
				claims: domain.Claims{
					UserID: noteEntity.ID,
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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(nil, fiber.NewError(fiber.StatusBadRequest, "validation error")).Once()
					return mockNoteService
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
			n := &NoteHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Post("/api/v1/note", n.Create)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/note", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.code, res.StatusCode)

			if tt.wantErr {
				var got domain.ErrorResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				t.Log(got)
			} else {
				var got domain.SuccessResponse
				err = json.NewDecoder(res.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Title, noteEntity.Title)
				assert.Equal(t, tt.args.req.Content, noteEntity.Content)
				assert.Equal(t, tt.args.req.Visibility, noteEntity.Visibility)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Title, got.Data.(map[string]interface{})["title"].(string))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Content, got.Data.(map[string]interface{})["content"].(string))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Visibility, got.Data.(map[string]interface{})["visibility"].(string))
			}
		})
	}
}

func TestNoteHandler_GetAll(t *testing.T) {
	type fields struct {
		service port.NoteService
	}

	noteEntity := []domain.Note{
		*noteEntity,
	}

	mockNoteService := mocks.NewNoteService(t)

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
				service: func() port.NoteService {
					mockNoteService.EXPECT().GetAll(mock.AnythingOfType("domain.NoteQuery"), mock.AnythingOfType("*domain.Metadata")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully retrieved notes",
				Data: domain.NotePaginationResponse{
					Notes: func() []domain.NoteResponse {
						var data []domain.NoteResponse
						for _, note := range noteEntity {
							data = append(data, domain.NoteResponse{
								ID:         note.ID,
								Title:      note.Title,
								Content:    note.Content,
								Visibility: note.Visibility,
								UserID:     note.UserID,
							})
						}
						return data
					}(),
					Metadata: metadataEntity,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().GetAll(mock.AnythingOfType("domain.NoteQuery"), mock.AnythingOfType("*domain.Metadata")).Return(nil, errors.New("failed")).Once()
					return mockNoteService
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
			n := &NoteHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Get("/api/v1/note", n.GetAll)

			req := httptest.NewRequest(fiber.MethodGet, "/api/v1/note", nil)
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

func TestNoteHandler_GetByID(t *testing.T) {
	type fields struct {
		service port.NoteService
	}

	type args struct {
		req    domain.NoteRequest
		claims *domain.Claims
	}

	mockNoteService := mocks.NewNoteService(t)

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
				service: func() port.NoteService {
					mockNoteService.EXPECT().GetByID(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("*domain.Claims")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
				claims: &domain.Claims{
					UserID: noteEntity.UserID,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully retrieved note by id",
				Data: domain.NoteResponse{
					ID:         noteEntity.ID,
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
					UserID:     noteEntity.UserID,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().GetByID(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("*domain.Claims")).Return(nil, errors.New("failed")).Once()
					return mockNoteService
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
				claims: &domain.Claims{
					UserID: noteEntity.UserID,
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
			n := &NoteHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", tt.args.claims)
				return ctx.Next()
			})
			app.Get("/api/v1/note/:id", n.GetByID)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/v1/note/%v", tt.args.req.ID), bytes.NewBuffer(requestBody))
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
				assert.Equal(t, tt.args.req.ID, noteEntity.ID)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Title, got.Data.(map[string]interface{})["title"].(string))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Content, got.Data.(map[string]interface{})["content"].(string))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Visibility, got.Data.(map[string]interface{})["visibility"].(string))
			}
		})
	}
}

func TestNoteHandler_Update(t *testing.T) {
	type fields struct {
		service   port.NoteService
		validator util.Validator
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
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Update(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("domain.Claims")).Return(noteEntity, nil).Once()
					return mockNoteService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.NoteRequest{
					ID:         noteEntity.ID,
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
				},
				claims: domain.Claims{
					UserID: noteEntity.UserID,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully updated note by id",
				Data: domain.NoteResponse{
					ID:         noteEntity.ID,
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
					UserID:     noteEntity.UserID,
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Update(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("domain.Claims")).Return(nil, errors.New("failed")).Once()
					return mockNoteService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.NoteRequest{
					ID:         noteEntity.ID,
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
				},
				claims: domain.Claims{
					UserID: noteEntity.ID,
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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Update(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("domain.Claims")).Return(nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockNoteService
				}(),
				validator: *validator,
			},
			args: args{
				req: domain.NoteRequest{
					ID:         noteEntity.ID,
					Title:      noteEntity.Title,
					Content:    noteEntity.Content,
					Visibility: noteEntity.Visibility,
				},
				claims: domain.Claims{
					UserID: noteEntity.ID,
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
				service:   mockNoteService,
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
			n := &NoteHandler{
				service:   tt.fields.service,
				validator: tt.fields.validator,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Put("/api/v1/note/:id", n.Update)

			requestBody, err := json.Marshal(tt.args.req)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/api/v1/note/%v", tt.args.req.ID), bytes.NewBuffer(requestBody))
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
				assert.Equal(t, tt.args.req.ID, noteEntity.ID)
				assert.Equal(t, tt.args.req.Title, noteEntity.Title)
				assert.Equal(t, tt.args.req.Content, noteEntity.Content)
				assert.Equal(t, tt.args.req.Visibility, noteEntity.Visibility)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).ID, uint(got.Data.(map[string]interface{})["id"].(float64)))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Title, got.Data.(map[string]interface{})["title"].(string))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Content, got.Data.(map[string]interface{})["content"].(string))
				assert.Equal(t, tt.want.(domain.SuccessResponse).Data.(domain.NoteResponse).Visibility, got.Data.(map[string]interface{})["visibility"].(string))
			}
		})
	}
}

func TestNoteHandler_Delete(t *testing.T) {
	type fields struct {
		service port.NoteService
	}

	type args struct {
		req    domain.NoteRequest
		claims domain.Claims
	}

	mockNoteService := mocks.NewNoteService(t)

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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Delete(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("domain.Claims")).Return(nil).Once()
					return mockNoteService
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
				claims: domain.Claims{
					UserID: noteEntity.UserID,
				},
			},
			code: fiber.StatusOK,
			want: domain.SuccessResponse{
				Message: "successfully deleted note by id",
			},
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				service: func() port.NoteService {
					mockNoteService.EXPECT().Delete(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("domain.Claims")).Return(errors.New("failed")).Once()
					return mockNoteService
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
				claims: domain.Claims{
					UserID: noteEntity.ID,
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
				service: func() port.NoteService {
					mockNoteService.EXPECT().Delete(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("domain.Claims")).Return(fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")).Once()
					return mockNoteService
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
				claims: domain.Claims{
					UserID: noteEntity.ID,
				},
			},
			code: fiber.StatusForbidden,
			want: domain.ErrorResponse{
				Message: "user does not have permission to perform this action",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NoteHandler{
				service: tt.fields.service,
			}

			app := config.NewFiber()
			app.Use(func(ctx *fiber.Ctx) error {
				ctx.Locals("claims", &tt.args.claims)
				return ctx.Next()
			})
			app.Delete("/api/v1/note/:id", n.Delete)

			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/v1/note/%v", tt.args.req.ID), nil)
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
				assert.Equal(t, tt.args.req.ID, noteEntity.ID)
				assert.Equal(t, tt.want.(domain.SuccessResponse).Message, got.Message)
			}
		})
	}
}
