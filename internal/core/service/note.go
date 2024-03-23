package service

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"

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

func (n *NoteService) GetAll() ([]domain.Note, error) {
	return n.repository.GetAll()
}

func (n *NoteService) GetByID(req domain.NoteRequest) (*domain.Note, error) {
	return n.repository.GetByID(req.ID)
}

func (n *NoteService) Update(req domain.NoteRequest, claims domain.Claims) (*domain.Note, error) {
	note, err := n.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if note.UserID != claims.UserID {
		return nil, fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return n.repository.Update(req, note)

}

func (n *NoteService) Delete(req domain.NoteRequest, claims domain.Claims) error {
	note, err := n.repository.GetByID(req.ID)
	if err != nil {
		return err
	}

	if note.UserID != claims.UserID {
		return fiber.NewError(fiber.StatusForbidden, "user does not have permission to perform this action")
	}

	return n.repository.Delete(note)
}
