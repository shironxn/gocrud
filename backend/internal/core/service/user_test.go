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

func TestUserService_Create(t *testing.T) {
	type fields struct {
		repository port.UserRepository
		bcrypt     util.Bcrypt
	}

	type args struct {
		req domain.UserRegisterRequest
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
					mockUserRepository.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     userEntity.Name,
					Email:    userEntity.Email,
					Password: "password123",
				},
			},
			want:    userEntity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}

			got, err := u.Create(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Name, userEntity.Name)
				assert.Equal(t, tt.args.req.Email, userEntity.Email)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	type fields struct {
		repository port.UserRepository
		bcrypt     util.Bcrypt
	}

	type args struct {
		req domain.UserLoginRequest
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
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    userEntity.Email,
					Password: "password123",
				},
			},
			want:    userEntity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("domain.UserRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "invalid password",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    userEntity.Email,
					Password: "invalid",
				},
			},
			want:    errors.New("invalid password"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}

			got, err := u.Login(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Email, userEntity.Email)
				assert.Equal(t, tt.want, got)
			}
		})
	}
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
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetAll(mock.AnythingOfType("domain.UserQuery"), mock.AnythingOfType("*domain.Metadata")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
			}

			got, err := u.GetAll(tt.args.req, &tt.args.metadata)

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
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Times(2).Once()
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
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(nil, errors.New("failed"))
					return mockUserRepository
				}(),
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
			}

			got, err := u.GetByID(tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, userEntity.ID)
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
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("*domain.User")).Return(userEntity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
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
			want:    userEntity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("*domain.User")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRequest{
					ID: userEntity.ID,
				},
				claims: domain.Claims{
					UserID: userEntity.ID,
				},
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
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
			u := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}

			got, err := u.Update(tt.args.req, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, userEntity.ID)
				assert.Equal(t, tt.args.req.Name, userEntity.Name)
				assert.Equal(t, tt.args.req.Email, userEntity.Email)
				assert.Equal(t, tt.args.req.Password, userEntity.Password)
				assert.Equal(t, tt.args.claims.UserID, userEntity.ID)
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
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
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
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
					mockUserRepository.EXPECT().Delete(mock.AnythingOfType("*domain.User")).Return(errors.New("failed")).Once()
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
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("domain.UserRequest")).Return(userEntity, nil).Once()
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
			u := &UserService{
				repository: tt.fields.repository,
			}

			err := u.Delete(tt.args.req, tt.args.claims)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.(error).Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.ID, userEntity.ID)
				assert.Equal(t, tt.args.claims.UserID, userEntity.ID)
			}
		})
	}
}