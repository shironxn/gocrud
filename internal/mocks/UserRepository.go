// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	domain "gocrud/internal/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

type UserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepository) EXPECT() *UserRepository_Expecter {
	return &UserRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: req
func (_m *UserRepository) Create(req domain.UserRegisterRequest) (*domain.User, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserRegisterRequest) (*domain.User, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(domain.UserRegisterRequest) *domain.User); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.UserRegisterRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - req domain.UserRegisterRequest
func (_e *UserRepository_Expecter) Create(req interface{}) *UserRepository_Create_Call {
	return &UserRepository_Create_Call{Call: _e.mock.On("Create", req)}
}

func (_c *UserRepository_Create_Call) Run(run func(req domain.UserRegisterRequest)) *UserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.UserRegisterRequest))
	})
	return _c
}

func (_c *UserRepository_Create_Call) Return(_a0 *domain.User, _a1 error) *UserRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Create_Call) RunAndReturn(run func(domain.UserRegisterRequest) (*domain.User, error)) *UserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: entity
func (_m *UserRepository) Delete(entity *domain.User) error {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - entity *domain.User
func (_e *UserRepository_Expecter) Delete(entity interface{}) *UserRepository_Delete_Call {
	return &UserRepository_Delete_Call{Call: _e.mock.On("Delete", entity)}
}

func (_c *UserRepository_Delete_Call) Run(run func(entity *domain.User)) *UserRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.User))
	})
	return _c
}

func (_c *UserRepository_Delete_Call) Return(_a0 error) *UserRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Delete_Call) RunAndReturn(run func(*domain.User) error) *UserRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields:
func (_m *UserRepository) GetAll() ([]domain.User, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type UserRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *UserRepository_Expecter) GetAll() *UserRepository_GetAll_Call {
	return &UserRepository_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *UserRepository_GetAll_Call) Run(run func()) *UserRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UserRepository_GetAll_Call) Return(_a0 []domain.User, _a1 error) *UserRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetAll_Call) RunAndReturn(run func() ([]domain.User, error)) *UserRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: email
func (_m *UserRepository) GetByEmail(email string) (*domain.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserRepository_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - email string
func (_e *UserRepository_Expecter) GetByEmail(email interface{}) *UserRepository_GetByEmail_Call {
	return &UserRepository_GetByEmail_Call{Call: _e.mock.On("GetByEmail", email)}
}

func (_c *UserRepository_GetByEmail_Call) Run(run func(email string)) *UserRepository_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UserRepository_GetByEmail_Call) Return(_a0 *domain.User, _a1 error) *UserRepository_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByEmail_Call) RunAndReturn(run func(string) (*domain.User, error)) *UserRepository_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: id
func (_m *UserRepository) GetByID(id uint) (*domain.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*domain.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *domain.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type UserRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - id uint
func (_e *UserRepository_Expecter) GetByID(id interface{}) *UserRepository_GetByID_Call {
	return &UserRepository_GetByID_Call{Call: _e.mock.On("GetByID", id)}
}

func (_c *UserRepository_GetByID_Call) Run(run func(id uint)) *UserRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *UserRepository_GetByID_Call) Return(_a0 *domain.User, _a1 error) *UserRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByID_Call) RunAndReturn(run func(uint) (*domain.User, error)) *UserRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: req, entity
func (_m *UserRepository) Update(req domain.UserRequest, entity *domain.User) (*domain.User, error) {
	ret := _m.Called(req, entity)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserRequest, *domain.User) (*domain.User, error)); ok {
		return rf(req, entity)
	}
	if rf, ok := ret.Get(0).(func(domain.UserRequest, *domain.User) *domain.User); ok {
		r0 = rf(req, entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.UserRequest, *domain.User) error); ok {
		r1 = rf(req, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type UserRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - req domain.UserRequest
//   - entity *domain.User
func (_e *UserRepository_Expecter) Update(req interface{}, entity interface{}) *UserRepository_Update_Call {
	return &UserRepository_Update_Call{Call: _e.mock.On("Update", req, entity)}
}

func (_c *UserRepository_Update_Call) Run(run func(req domain.UserRequest, entity *domain.User)) *UserRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.UserRequest), args[1].(*domain.User))
	})
	return _c
}

func (_c *UserRepository_Update_Call) Return(_a0 *domain.User, _a1 error) *UserRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Update_Call) RunAndReturn(run func(domain.UserRequest, *domain.User) (*domain.User, error)) *UserRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
