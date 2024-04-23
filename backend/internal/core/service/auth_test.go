package service

import (
	"errors"
	"testing"

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

func TestUserService_Register(t *testing.T) {
	type fields struct {
		repository port.AuthRepository
		bcrypt     util.Bcrypt
	}

	type args struct {
		req domain.AuthRegisterRequest
	}

	mockAuthRepository := mocks.NewAuthRepository(t)
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
				repository: func() port.AuthRepository {
					mockAuthRepository.EXPECT().Register(mock.AnythingOfType("domain.AuthRegisterRequest")).Return(authEntity, nil).Once()
					return mockAuthRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.AuthRegisterRequest{
					Name:     authEntity.Name,
					Email:    authEntity.Email,
					Password: "password123",
				},
			},
			want:    authEntity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.AuthRepository {
					mockAuthRepository.EXPECT().Register(mock.AnythingOfType("domain.AuthRegisterRequest")).Return(nil, errors.New("failed")).Once()
					return mockAuthRepository
				}(),
				bcrypt: bcrypt,
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}

			got, err := h.Register(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Name, authEntity.Name)
				assert.Equal(t, tt.args.req.Email, authEntity.Email)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	type fields struct {
		repository port.AuthRepository
		bcrypt     util.Bcrypt
		jwt        util.JWT
	}

	type args struct {
		req domain.AuthLoginRequest
	}

	mockAuthRepository := mocks.NewAuthRepository(t)
	bcrypt := util.NewBcrypt()
	jwt := util.NewJWT(&config.Config{})

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
				repository: func() port.AuthRepository {
					mockAuthRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(authEntity, nil).Once()
					mockAuthRepository.EXPECT().StoreRefreshToken(mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil).Once()
					return mockAuthRepository
				}(),
				bcrypt: bcrypt,
				jwt:    jwt,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: "password123",
				},
			},
			want:    authEntity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.AuthRepository {
					mockAuthRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(nil, errors.New("failed")).Once()
					return mockAuthRepository
				}(),
				bcrypt: bcrypt,
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "invalid password",
			fields: fields{
				repository: func() port.AuthRepository {
					mockAuthRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(authEntity, nil).Once()
					return mockAuthRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.AuthLoginRequest{
					Email:    authEntity.Email,
					Password: "invalid",
				},
			},
			want:    errors.New("invalid password"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
				jwt:        tt.fields.jwt,
			}

			got, tokens, err := h.Login(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.Equal(t, tt.args.req.Email, authEntity.Email)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
