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

func (n *NoteService) Create(req domain.NoteRequest) (*domain.Note, error) {
	return n.repository.Create(req)
}

func (n *NoteService) GetAll(req domain.NoteQuery, metadata *domain.Metadata) ([]domain.Note, error) {
	data, err := n.repository.GetAll(req, metadata)
	return data, err
}

func (n *NoteService) GetByID(req domain.NoteRequest, claims *domain.Claims) (*domain.Note, error) {
	data, err := n.repository.GetByID(req)
	if err != nil {
		return nil, err
	}

	if data.Visibility == "private" && (claims == nil || data.UserID != claims.UserID) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "you are not authorized to access this private note")
	}

	return data, nil
}

func (n *NoteService) Update(req domain.NoteRequest, claims domain.Claims) (*domain.Note, error) {
	note, err := n.repository.GetByID(req)
	if err != nil {
		return nil, err
	}

	if note.UserID != claims.UserID {
		return nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return n.repository.Update(req, note)

}

func (n *NoteService) Delete(req domain.NoteRequest, claims domain.Claims) error {
	note, err := n.repository.GetByID(req)
	if err != nil {
		return err
	}

	if note.UserID != claims.UserID {
		return fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return n.repository.Delete(note)
}
