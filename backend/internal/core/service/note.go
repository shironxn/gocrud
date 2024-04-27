package service

import (
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type NoteService struct {
	repository port.NoteRepository
}

func NewNoteService(repository port.NoteRepository) port.NoteService {
	return &NoteService{
		repository: repository,
	}
}

func (h *NoteService) Create(req domain.NoteRequest) (*domain.Note, error) {
	return h.repository.Create(req)
}

func (h *NoteService) GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error) {
	data, err := h.repository.GetAll(req, metadata)
	return data, err
}

func (h *NoteService) GetByID(id uint, claims *domain.Claims) (*domain.Note, error) {
	data, err := h.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if data.Visibility == "private" && (claims == nil || data.UserID != claims.UserID) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "you are not authorized to access this private note")
	}

	return data, nil
}

func (h *NoteService) Update(req domain.NoteUpdate, claims domain.Claims) (*domain.Note, error) {
	note, err := h.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if note.UserID != claims.UserID {
		return nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return h.repository.Update(req, note)
}

func (h *NoteService) Delete(id uint, claims domain.Claims) error {
	note, err := h.repository.GetByID(id)
	if err != nil {
		return err
	}

	if note.UserID != claims.UserID {
		return fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return h.repository.Delete(note)
}
