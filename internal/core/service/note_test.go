package service

import (
	"errors"
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"
	"gocrud/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNoteService_Create(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req domain.NoteRequest
	}

	entity := &domain.Note{
		Model: gorm.Model{
			ID: 1,
		},
		Title:      "golang",
		Content:    "is the best",
		Visibility: "public",
		UserID:     1,
	}

	mockNoteRepository := mocks.NewNoteRepository(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(entity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					Title:      entity.Title,
					Content:    entity.Content,
					Visibility: entity.Visibility,
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(nil, errors.New("failed")).Once()
					return mockNoteRepository
				}(),
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := u.Create(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Title, entity.Title)
				assert.Equal(t, tt.args.req.Content, entity.Content)
				assert.Equal(t, tt.args.req.Visibility, entity.Visibility)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNoteService_GetAll(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	entity := []domain.Note{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Title:      "golang",
			Content:    "is the best",
			Visibility: "public",
			UserID:     1,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Title:      "golang",
			Content:    "is the best",
			Visibility: "public",
			UserID:     1,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Title:      "golang",
			Content:    "is the best",
			Visibility: "public",
			UserID:     1,
		},
	}

	mockNoteRepository := mocks.NewNoteRepository(t)

	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetAll().Return(entity, nil).Once()
					return mockNoteRepository
				}(),
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetAll().Return(nil, errors.New("failed")).Once()
					return mockNoteRepository
				}(),
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := u.GetAll()

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNoteService_GetByID(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req domain.NoteRequest
	}

	entity := &domain.Note{
		Model: gorm.Model{
			ID: 1,
		},
		Title:      "golang",
		Content:    "is the best",
		Visibility: "public",
		UserID:     1,
	}

	mockNoteRepository := mocks.NewNoteRepository(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: entity.ID,
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(nil, errors.New("failed")).Once()
					return mockNoteRepository
				}(),
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := u.GetByID(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, entity.ID)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNoteService_Update(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req    domain.NoteRequest
		claims domain.Claims
	}

	entity := &domain.Note{
		Model: gorm.Model{
			ID: 1,
		},
		Title:      "golang",
		Content:    "is the best",
		Visibility: "public",
		UserID:     1,
	}

	mockNoteRepository := mocks.NewNoteRepository(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockNoteRepository.EXPECT().Update(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("*domain.Note")).Return(entity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID:         entity.ID,
					Title:      entity.Title,
					Content:    entity.Content,
					Visibility: entity.Visibility,
				},
				claims: domain.Claims{
					UserID: entity.UserID,
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockNoteRepository.EXPECT().Update(mock.AnythingOfType("domain.NoteRequest"), mock.AnythingOfType("*domain.Note")).Return(nil, errors.New("failed")).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.UserID,
				},
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.UserID + 1,
				},
			},
			want:    errors.New("user does not have permission to perform this action"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := u.Update(tt.args.req, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, entity.ID)
				assert.Equal(t, tt.args.req.Title, entity.Title)
				assert.Equal(t, tt.args.req.Content, entity.Content)
				assert.Equal(t, tt.args.req.Visibility, entity.Visibility)
				assert.Equal(t, tt.args.claims.UserID, entity.UserID)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNoteService_Delete(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req    domain.NoteRequest
		claims domain.Claims
	}

	entity := &domain.Note{
		Model: gorm.Model{
			ID: 1,
		},
		Title:      "golang",
		Content:    "is the best",
		Visibility: "public",
		UserID:     1,
	}

	mockNoteRepository := mocks.NewNoteRepository(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockNoteRepository.EXPECT().Delete(mock.AnythingOfType("*domain.Note")).Return(nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.ID,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockNoteRepository.EXPECT().Delete(mock.AnythingOfType("*domain.Note")).Return(errors.New("failed")).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.UserID,
				},
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.UserID + 1,
				},
			},
			want:    errors.New("user does not have permission to perform this action"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &NoteService{
				repository: tt.fields.repository,
			}

			err := u.Delete(tt.args.req, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, entity.ID)
				assert.Equal(t, tt.args.claims.UserID, entity.UserID)
				assert.Equal(t, tt.args.claims.UserID, entity.UserID)
			}
		})
	}
}
