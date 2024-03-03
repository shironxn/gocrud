package port

import (
	"gocrud/internal/core/domain"
	"net/http"
)

type UserRepository interface {
	Create(req domain.UserRequest) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(entity *domain.User, req domain.UserRequest) (*domain.User, error)
	Delete(entity *domain.User) error
}

type UserService interface {
	Create(req domain.UserRequest) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(req domain.UserRequest) (*domain.User, error)
	Update(req domain.UserRequest, claims domain.Claims) (*domain.User, error)
	Delete(req domain.UserRequest, claims domain.Claims) error
}

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
