package service

import (
	"errors"
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"
	"gocrud/internal/mocks"
	"gocrud/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserService_Create(t *testing.T) {
	type fields struct {
		repository port.UserRepository
		bcrypt     util.Bcrypt
	}

	type args struct {
		req domain.UserRegisterRequest
	}

	entity := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserRepository.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     entity.Name,
					Email:    entity.Email,
					Password: "password123",
				},
			},
			want:    entity,
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
				assert.Equal(t, tt.args.req.Name, entity.Name)
				assert.Equal(t, tt.args.req.Email, entity.Email)
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

	entity := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    entity.Email,
					Password: "password123",
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(nil, errors.New("failed")).Once()
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
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    entity.Email,
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
				assert.Equal(t, tt.args.req.Email, entity.Email)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}

	entity := []domain.User{
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

	mockUserRepository := mocks.NewUserRepository(t)

	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetAll().Return(entity, nil).Once()
					return mockUserRepository
				}(),
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetAll().Return(nil, errors.New("failed")).Once()
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

func TestUserService_GetByID(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}

	type args struct {
		req domain.UserRequest
	}

	entity := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: entity.ID,
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(nil, errors.New("failed")).Once()
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
				assert.Equal(t, tt.args.req.ID, entity.ID)
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

	entity := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("*domain.User")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRequest{
					ID:       entity.ID,
					Name:     entity.Name,
					Email:    entity.Email,
					Password: entity.Password,
				},
				claims: domain.Claims{
					UserID: entity.ID,
				},
			},
			want:    entity,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("domain.UserRequest"), mock.AnythingOfType("*domain.User")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.ID,
				},
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.ID + 1,
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
				assert.Equal(t, tt.args.req.ID, entity.ID)
				assert.Equal(t, tt.args.req.Name, entity.Name)
				assert.Equal(t, tt.args.req.Email, entity.Email)
				assert.Equal(t, tt.args.req.Password, entity.Password)
				assert.Equal(t, tt.args.claims.UserID, entity.ID)
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

	entity := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
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
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockUserRepository.EXPECT().Delete(mock.AnythingOfType("*domain.User")).Return(nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
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
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					mockUserRepository.EXPECT().Delete(mock.AnythingOfType("*domain.User")).Return(errors.New("failed")).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.ID,
				},
			},
			want:    errors.New("failed"),
			wantErr: true,
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(entity, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: entity.ID,
				},
				claims: domain.Claims{
					UserID: entity.ID + 1,
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
				assert.Equal(t, tt.args.req.ID, entity.ID)
				assert.Equal(t, tt.args.claims.UserID, entity.ID)
			}
		})
	}
}
