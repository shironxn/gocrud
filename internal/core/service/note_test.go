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
					Title:      "golang",
					Content:    "is the best",
					Visibility: "public",
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				repository: func() port.NoteRepository {
					mockNoteRepository.EXPECT().Create(mock.AnythingOfType("domain.NoteRequest")).Return(nil, errors.New("failed")).Once()
					return mockNoteRepository
				}(),
			},
			args: args{
				req: domain.NoteRequest{
					Title:      "golang",
					Content:    "is the best",
					Visibility: "public",
				},
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
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
