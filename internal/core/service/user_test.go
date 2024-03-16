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
		bcrypt     *util.Bcrypt
	}
	type args struct {
		req domain.UserRegisterRequest
	}
	expected := &domain.User{
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
		want    *domain.User
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().Create(mock.AnythingOfType("domain.UserRegisterRequest")).Return(expected, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRegisterRequest{
					Name:     "shiron",
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			want: &domain.User{
				Model:    gorm.Model{ID: 1},
				Name:     "shiron",
				Email:    "shiron@example.com",
				Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
			},
			wantErr: nil,
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
			args: args{
				req: domain.UserRegisterRequest{
					Name:     "shiron",
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			want:    nil,
			wantErr: errors.New("failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}
			got, err := u.Create(tt.args.req)
			if (err != nil) && (tt.wantErr == nil) {
				t.Errorf("UserService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	type fields struct {
		repository port.UserRepository
		bcrypt     *util.Bcrypt
	}
	type args struct {
		req domain.UserLoginRequest
	}
	expected := &domain.User{
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
		want    *domain.User
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByEmail(mock.AnythingOfType("string")).Return(expected, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserLoginRequest{
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			want:    expected,
			wantErr: nil,
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
			args: args{
				req: domain.UserLoginRequest{
					Email:    "shiron@example.com",
					Password: "password123",
				},
			},
			want:    nil,
			wantErr: errors.New("failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}
			got, err := u.Login(tt.args.req)
			if (err != nil) && (tt.wantErr == nil) {
				t.Errorf("UserService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}
	expected := []domain.User{
		{
			Model:    gorm.Model{ID: 1},
			Name:     "shiron",
			Email:    "shiron@example.com",
			Password: "$2y$10$YovD7LTJb0XqE.Ll1Xtjnuns6tHiQM7MdO5T2QuThx3UyfLCkP1o6",
		},
	}
	mockUserRepository := mocks.NewUserRepository(t)
	tests := []struct {
		name    string
		fields  fields
		want    []domain.User
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetAll().Return(expected, nil).Once()
					return mockUserRepository
				}(),
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetAll().Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
			},
			want:    nil,
			wantErr: errors.New("failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
			}
			got, err := u.GetAll()
			if (err != nil) && (tt.wantErr == nil) {
				t.Errorf("UserService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
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
	expected := &domain.User{
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
		want    *domain.User
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: 1,
				},
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
			},
			want:    nil,
			wantErr: errors.New("failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
			}
			got, err := u.GetByID(tt.args.req)
			if (err != nil) && (tt.wantErr == nil) {
				t.Errorf("UserService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	type fields struct {
		repository port.UserRepository
		bcrypt     *util.Bcrypt
	}
	type args struct {
		req    domain.UserRequest
		claims *domain.Claims
	}
	expected := &domain.User{
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
		want    *domain.User
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("*domain.User"), mock.AnythingOfType("domain.UserRequest")).Return(expected, nil).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRequest{
					Name:     "shiron",
					Email:    "shiron@example.com",
					Password: "password123",
				},
				claims: &domain.Claims{
					UserID: 1,
				},
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					mockUserRepository.EXPECT().Update(mock.AnythingOfType("*domain.User"), mock.AnythingOfType("domain.UserRequest")).Return(nil, errors.New("failed")).Once()
					return mockUserRepository
				}(),
				bcrypt: bcrypt,
			},
			args: args{
				req: domain.UserRequest{
					Name:     "shiron",
					Email:    "shiron@example.com",
					Password: "password123",
				},
				claims: &domain.Claims{
					UserID: 1,
				},
			},
			want:    nil,
			wantErr: errors.New("failed"),
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: 1,
				},
				claims: &domain.Claims{
					UserID: 2,
				},
			},
			wantErr: errors.New("user does not have permission to perform this action"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
				bcrypt:     tt.fields.bcrypt,
			}
			got, err := u.Update(tt.args.req, tt.args.claims)
			if (err != nil) && (tt.wantErr == nil) {
				t.Errorf("UserService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	type fields struct {
		repository port.UserRepository
	}
	type args struct {
		req    domain.UserRequest
		claims *domain.Claims
	}
	expected := &domain.User{
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
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					mockUserRepository.EXPECT().Delete(mock.AnythingOfType("*domain.User")).Return(nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: 10,
				},
				claims: &domain.Claims{
					UserID: 1,
				},
			},
			wantErr: nil,
		},
		{
			name: "failure",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					mockUserRepository.EXPECT().Delete(mock.AnythingOfType("*domain.User")).Return(errors.New("failed")).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: 1,
				},
				claims: &domain.Claims{
					UserID: 1,
				},
			},
			wantErr: errors.New("failed"),
		},
		{
			name: "permission denied",
			fields: fields{
				repository: func() port.UserRepository {
					mockUserRepository.EXPECT().GetByID(mock.AnythingOfType("uint")).Return(expected, nil).Once()
					return mockUserRepository
				}(),
			},
			args: args{
				req: domain.UserRequest{
					ID: 1,
				},
				claims: &domain.Claims{
					UserID: 2,
				},
			},
			wantErr: errors.New("user does not have permission to perform this action"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				repository: tt.fields.repository,
			}
			err := u.Delete(tt.args.req, tt.args.claims)
			if (err != nil) && (tt.wantErr == nil) {
				t.Errorf("UserService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
