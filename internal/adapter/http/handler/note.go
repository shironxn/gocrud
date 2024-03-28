package handler

import (
	"gocrud/internal/core/domain"
	"gocrud/internal/core/port"
	"gocrud/internal/util"

	_ "gocrud/docs"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	service   port.NoteService
	validator util.Validator
}

func NewNoteHandler(service port.NoteService, validator util.Validator) port.NoteHandler {
	return &NoteHandler{
		service:   service,
		validator: validator,
	}
}

// @Summary Create a new note
// @Description Create a new note with the specified title, content, and visibility
// @Tags note
// @Accept json
// @Produce json
// @Param note body domain.NoteRequest true "Note request object"
// @Success 201 {object} domain.NoteResponse "Successfully created a new note"
// @Failure 400 {object} domain.ErrorValidationResponse "Validation error"
// @Router /note [post]
func (n *NoteHandler) Create(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims := ctx.Locals("claims").(*domain.Claims)
	req.UserID = claims.UserID

	if err := n.validator.Validate(req); err != nil {
		if err := n.validator.Validate(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err)
		}
	}

	result, err := n.service.Create(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		domain.SuccessResponse{
			Message: "successfully created note",
			Data: domain.NoteResponse{
				ID:         result.ID,
				Title:      result.Title,
				Content:    result.Content,
				Visibility: result.Visibility,
				UserID:     result.UserID,
				UpdatedAt:  result.UpdatedAt,
				CreatedAt:  result.CreatedAt,
			},
		},
	)
}

// @Summary Get all notes
// @Description Retrieve all available notes
// @Tags note
// @Produce json
// @Success 200 {object} []domain.NoteResponse "Successfully retrieved all notes"
// @Router /note [get]
func (n *NoteHandler) GetAll(ctx *fiber.Ctx) error {
	result, err := n.service.GetAll()
	if err != nil {
		return err
	}

	var data []domain.NoteResponse
	for _, note := range result {
		data = append(data, domain.NoteResponse{
			ID:         note.ID,
			Title:      note.Title,
			Content:    note.Content,
			Visibility: note.Visibility,
			UserID:     note.UserID,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully retrieved notes",
		Data:    data,
	})
}

// @Summary Get a note by ID
// @Description Retrieve a note based on the provided ID
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} domain.NoteResponse "Successfully retrieved a note by ID"
// @Router /note/{id} [get]
func (n *NoteHandler) GetByID(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	result, err := n.service.GetByID(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully retrieved note by id",
		Data: domain.NoteResponse{
			ID:         result.ID,
			Title:      result.Title,
			Content:    result.Content,
			Visibility: result.Visibility,
			UserID:     result.UserID,
			UpdatedAt:  result.UpdatedAt,
			CreatedAt:  result.CreatedAt,
		},
	})
}

// @Summary Update a note by ID
// @Description Update an existing note based on the provided ID
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param note body domain.NoteRequest true "Updated note object"
// @Success 200 {object} domain.NoteResponse "Successfully updated a note by ID"
// @Router /note/{id} [put]
func (n *NoteHandler) Update(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims := ctx.Locals("claims").(*domain.Claims)
	req.UserID = claims.UserID

	if err := n.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := n.service.Update(req, *claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully updated note by id",
		Data: domain.NoteResponse{
			ID:         result.ID,
			Title:      result.Title,
			Content:    result.Content,
			Visibility: result.Visibility,
			UserID:     result.UserID,
			UpdatedAt:  result.UpdatedAt,
			CreatedAt:  result.CreatedAt,
		},
	})
}

// @Summary Delete a note by ID
// @Description Delete an existing note based on the provided ID
// @Tags note
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} domain.SuccessResponse "Successfully deleted a note by ID"
// @Router /note/{id} [delete]
func (n *NoteHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	claims := ctx.Locals("claims").(*domain.Claims)

	if err := n.service.Delete(req, *claims); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully deleted note by id",
	})
}
