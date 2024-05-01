package service

import (
	"errors"
	"testing"

	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/mocks"
	"github.com/shironxn/gocrud/internal/util"

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

func TestUserService_GetAll(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}

	type args struct {
		req      domain.UserQuery
		metadata domain.Metadata
	}

	userEntity := []domain.User{
		*userEntity,
	}

	mockUserRepository := mocks.NewUserRepository(t)

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
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetAll(mock.AnythingOfType("domain.UserQuery"), mock.AnythingOfType("*domain.Metadata")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req:      domain.UserQuery{},
				metadata: domain.Metadata{},
			},
			want:    userEntity,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserService{
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

func TestUserService_GetByID(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}

	type args struct {
		req domain.UserRequest
	}

	mockUserRepository := mocks.NewUserRepository(t)

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
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(userEntity, nil).Times(2).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
			},
			want:    userEntity,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserService{
				repository: tt.fields.repository,
			}

			got, err := h.GetByID(tt.args.req.ID)

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

func TestUserService_Update(t *testing.T) {
	type fields struct {
		repository port.UserRepository
		bcrypt     util.Bcrypt
	}

	type args struct {
		req    domain.UserRequest
		claims domain.Claims
	}

	mockUserRepository := mocks.NewUserRepository(t)
	bcrypt := util.NewBcrypt()

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
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(userEntity, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("*domain.User")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRequest{},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			want:    userEntity,
			wantErr: false,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(userEntity.ID).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
				claims: domain.Claims{
					UserID: userEntity.ID + 1,
				},
			},
			want:    errors.New("user does not have permission to perform this action"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
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

func TestUserService_Delete(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}

	type args struct {
		req    domain.UserRequest
		claims domain.Claims
	}

	mockUserRepository := mocks.NewUserRepository(t)

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
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(userEntity, nil).Once()
					mockUserRepository.EXPECT().Delete(mock.AnythingOfType("*domain.User")).Return(nil).Once()
					return mockUserRepository
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
			want:    nil,
			wantErr: false,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
				claims: domain.Claims{
					UserID: userEntity.ID + 1,
				},
			},
			want:    errors.New("user does not have permission to perform this action"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserService{
				repository: tt.fields.repository,
			}

			err := h.Delete(tt.args.req, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
