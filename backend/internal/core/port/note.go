package port

import (
	"github.com/shironxn/gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type NoteRepository interface {
	Create(req domain.NoteRequest) (*domain.Note, error)
	GetAll(req domain.NoteQuery, metdata *domain.Metadata) ([]domain.Note, error)
	GetByID(req domain.NoteRequest) (*domain.Note, error)
	Update(req domain.NoteRequest, entity *domain.Note) (*domain.Note, error)
	Delete(entity *domain.Note) error
}

type NoteService interface {
	Create(req domain.NoteRequest) (*domain.Note, error)
	GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error)
	GetByID(req domain.NoteRequest, claims *domain.Claims) (*domain.Note, error)
	Update(req domain.NoteRequest, claims domain.Claims) (*domain.Note, error)
	Delete(req domain.NoteRequest, claims domain.Claims) error
}

type NoteHandler interface {
	Create(ctx *fiber.Ctx) error
	GetAll(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
