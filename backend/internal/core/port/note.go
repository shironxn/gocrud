package port

import (
	"github.com/shironxn/gocrud/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type NoteRepository interface {
	Create(req domain.NoteRequest) (*domain.Note, error)
	GetAll(req domain.NoteQuery, metdata *domain.Metadata) ([]domain.Note, error)
	GetByID(id uint) (*domain.Note, error)
	Update(req domain.NoteUpdateRequest, note *domain.Note) (*domain.Note, error)
	Delete(note *domain.Note) error
}

type NoteService interface {
	Create(req domain.NoteRequest) (*domain.Note, error)
	GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error)
	GetByID(id uint, claims *domain.Claims) (*domain.Note, error)
	Update(req domain.NoteUpdateRequest, claims domain.Claims) (*domain.Note, error)
	Delete(id uint, claims domain.Claims) error
}

type NoteHandler interface {
	Create(ctx *fiber.Ctx) error
	GetAll(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
