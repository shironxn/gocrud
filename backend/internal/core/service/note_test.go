package service

import (
	"errors"
	"testing"

	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/mocks"

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

func TestNoteService_Create(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req domain.NoteRequest
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
					mockNoteRepository.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(noteEntity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{},
			},
			want:    noteEntity,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := h.Create(tt.args.req)

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

func TestNoteService_GetAll(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req      domain.NoteQuery
		metadata domain.Metadata
	}

	noteEntity := []domain.Note{
		*noteEntity,
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
					mockNoteRepository.EXPECT().GetAll(mock.AnythingOfType("domain.NoteQuery"), mock.AnythingOfType("*domain.Metadata")).Return(noteEntity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req:      domain.NoteQuery{},
				metadata: domain.Metadata{},
			},
			want:    noteEntity,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := h.GetAll(tt.args.req, &tt.args.metadata)

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
		req    domain.NoteRequest
		claims *domain.Claims
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
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(noteEntity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
			},
			want:    noteEntity,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := h.GetByID(tt.args.req.ID, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, noteEntity.ID)
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
		req    domain.NoteUpdateRequest
		claims domain.Claims
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
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(noteEntity, nil).Once()
					mockNoteRepository.EXPECT().Update(mock.AnythingOfType("domain.NoteUpdateRequest"), mock.AnythingOfType("*domain.Note")).Return(noteEntity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteUpdateRequest{},
				claims: domain.Claims{
					UserID: noteEntity.UserID,
				},
			},
			want:    noteEntity,
			wantErr: false,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(noteEntity.UserID).Return(noteEntity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteUpdateRequest{
					ID: noteEntity.ID,
				},
				claims: domain.Claims{
					UserID: noteEntity.UserID + 1,
				},
			},
			want:    errors.New("user does not have permission to perform this action"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteService{
				repository: tt.fields.repository,
			}

			got, err := h.Update(tt.args.req, tt.args.claims)

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

func TestNoteService_Delete(t *testing.T) {
	type fields struct {
		repository port.NoteRepository
	}

	type args struct {
		req    domain.NoteRequest
		claims domain.Claims
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
					mockNoteRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(noteEntity, nil).Once()
					mockNoteRepository.EXPECT().Delete(mock.AnythingOfType("*domain.Note")).Return(nil).Once()
					return mockNoteRepository
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
			want:    nil,
			wantErr: false,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().GetByID(noteEntity.ID).Return(noteEntity, nil).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					ID: noteEntity.ID,
				},
				claims: domain.Claims{
					UserID: noteEntity.UserID + 1,
				},
			},
			want:    errors.New("user does not have permission to perform this action"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteService{
				repository: tt.fields.repository,
			}

			err := h.Delete(tt.args.req.ID, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
