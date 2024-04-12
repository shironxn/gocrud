package handler

import (
	"reflect"
	"strconv"

	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	_ "github.com/shironxn/gocrud/docs"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	service   port.NoteService
	validator util.Validator
	jwt       util.JWT
}

func NewNoteHandler(service port.NoteService, validator util.Validator, jwt util.JWT) port.NoteHandler {
	return &NoteHandler{
		service:   service,
		validator: validator,
		jwt:       jwt,
	}
}

// @Summary Create a new note
// @Description Create a new note with the specified title, content, and visibility
// @Tags note
// @Accept json
// @Produce json
// @Param note body domain.NoteRequest true "Note request object"
// @Success 201 {object} domain.NoteResponse "Successfully created a new note"
// @Router /notes [post]
func (n *NoteHandler) Create(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := n.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
	req.UserID = claims.UserID

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
// @Param title query string false "Filter notes by title"
// @Param user_id query string false "Filter notes by user ID"
// @Param visibility query string false "Filter notes by visibility"
// @Success 200 {object} domain.NotePaginationResponse "Successfully retrieved all notes"
// @Router /notes [get]
func (n *NoteHandler) GetAll(ctx *fiber.Ctx) error {
	var req domain.NoteQuery
	var metadata domain.Metadata
	var data []domain.NoteResponse

	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

	if err := ctx.QueryParser(&metadata); err != nil {
		return err
	}

	switch {
	case req.Visibility == "private" || req.UserID != "" || req.Title != "":
		cookie := ctx.Cookies("token")
		claims, err := n.jwt.ValidateToken(cookie)
		if err != nil {
			req.Visibility = "public"
			break
		}
		if req.UserID == claims.ID {
			req.UserID = strconv.Itoa(int(claims.UserID))
		}
	default:
		req.Visibility = "public"
	}

	result, err := n.service.GetAll(req, &metadata)
	if err != nil {
		return err
	}

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
		Data: domain.NotePaginationResponse{
			Notes:    data,
			Metadata: metadata,
		},
	})
}

// @Summary Get a note by ID
// @Description Retrieve a note based on the provided ID
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} domain.NoteResponse "Successfully retrieved a note by ID"
// @Router /notes/{id} [get]
func (n *NoteHandler) GetByID(ctx *fiber.Ctx) error {
	var req domain.NoteRequest
	var result *domain.Note

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	cookie := ctx.Cookies("token")
	claims, err := n.jwt.ValidateToken(cookie)
	if err != nil {
		data, err := n.service.GetByID(req, nil)
		if err != nil {
			return err
		}
		result = data
	} else {
		data, err := n.service.GetByID(req, claims)
		if err != nil {
			return err
		}
		result = data
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
// @Router /notes/{id} [put]
func (n *NoteHandler) Update(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if reflect.DeepEqual(req, domain.NoteRequest{}) {
		return fiber.NewError(fiber.StatusBadRequest, "at least one field must be filled")
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
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
// @Router /notes/{id} [delete]
func (n *NoteHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
	if err := n.service.Delete(req, *claims); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.SuccessResponse{
		Message: "successfully deleted note by id",
	})
}
