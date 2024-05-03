package handler

import (
	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/port"
	"github.com/shironxn/gocrud/internal/util"

	_ "github.com/shironxn/gocrud/docs"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	service   port.NoteService
	validator *util.Validator
	jwt       util.JWT
	cfg       *config.Config
}

func NewNoteHandler(service port.NoteService, validator *util.Validator, jwt util.JWT, cfg *config.Config) port.NoteHandler {
	return &NoteHandler{
		service:   service,
		validator: validator,
		jwt:       jwt,
		cfg:       cfg,
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
func (h *NoteHandler) Create(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
	req.UserID = claims.UserID

	result, err := h.service.Create(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(domain.NoteResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		CoverURL:    result.CoverURL,
		Content:     result.Content,
		Visibility:  string(result.Visibility),
		Author: domain.NoteAuthor{
			ID:        result.Author.ID,
			Name:      result.Author.Name,
			AvatarURL: result.Author.AvatarURL,
		},
		UpdatedAt: result.UpdatedAt,
		CreatedAt: result.CreatedAt,
	})
}

// @Summary Get all notes
// @Description Retrieve all available notes
// @Tags note
// @Produce json
// @Param title query string false "Filter notes by title"
// @Param author query string false "Filter notes by author"
// @Param user_id query string false "Filter notes by user ID"
// @Param visibility query string false "Filter notes by visibility"
// @Param sort query string false "Sorting (e.g., +title, -created_at)"
// @Param order query string false "Sort order (e.g., asc, desc)"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} domain.NotePaginationResponse "Successfully retrieved all notes"
// @Router /notes [get]
func (h *NoteHandler) GetAll(ctx *fiber.Ctx) error {
	var req domain.NoteQuery
	var metadata domain.Metadata
	var data []domain.NoteResponse

	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

	if err := ctx.QueryParser(&metadata); err != nil {
		return err
	}

	req.Visibility = "public"
	cookie := ctx.Cookies("access-token")
	if cookie != "" {
		claims, _ := h.jwt.ValidateToken(cookie, h.cfg.JWT.Access)
		if claims != nil {
			switch {
			case req.Visibility == "private":
				req.UserID = int(claims.UserID)
			case req.UserID != 0:
				req.UserID = int(claims.UserID)
			default:
				req.Visibility = "public"
			}
		} else {
			req.Visibility = "public"
		}
	}

	result, err := h.service.GetAll(req, &metadata)
	if err != nil {
		return err
	}

	for _, note := range result {
		data = append(data, domain.NoteResponse{
			ID:          note.ID,
			Title:       note.Title,
			Description: note.Description,
			CoverURL:    note.CoverURL,
			Content:     note.Content,
			Visibility:  string(note.Visibility),
			Author: domain.NoteAuthor{
				ID:        note.Author.ID,
				Name:      note.Author.Name,
				AvatarURL: note.Author.AvatarURL,
			},
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.NotePaginationResponse{
		Notes:    data,
		Metadata: metadata,
	})
}

// @Summary Get a note by ID
// @Description Retrieve a note based on the provided ID
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} domain.NoteResponse "Successfully retrieved a note by ID"
// @Router /notes/{id} [get]
func (h *NoteHandler) GetByID(ctx *fiber.Ctx) error {
	var req domain.NoteRequest
	var result *domain.Note

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	cookie := ctx.Cookies("access-token")
	if cookie != "" {
		claims, err := h.jwt.ValidateToken(cookie, h.cfg.JWT.Access)
		if err != nil {
			data, err := h.service.GetByID(req.ID, nil)
			if err != nil {
				return err
			}
			result = data
		} else {
			data, err := h.service.GetByID(req.ID, claims)
			if err != nil {
				return err
			}
			result = data
		}
	} else {
		data, err := h.service.GetByID(req.ID, nil)
		if err != nil {
			return err
		}
		result = data
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.NoteResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		CoverURL:    result.CoverURL,
		Content:     result.Content,
		Visibility:  string(result.Visibility),
		Author: domain.NoteAuthor{
			ID:        result.Author.ID,
			Name:      result.Author.Name,
			AvatarURL: result.Author.AvatarURL,
		},
		UpdatedAt: result.UpdatedAt,
		CreatedAt: result.CreatedAt,
	})
}

// @Summary Update a note by ID
// @Description Update an existing note based on the provided ID
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param note body domain.NoteUpdateRequest true "Updated note object"
// @Success 200 {object} domain.NoteResponse "Successfully updated a note by ID"
// @Router /notes/{id} [put]
func (h *NoteHandler) Update(ctx *fiber.Ctx) error {
	var req domain.NoteUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
	req.UserID = claims.UserID

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := h.service.Update(req, *claims)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(domain.NoteResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		CoverURL:    result.CoverURL,
		Content:     result.Content,
		Visibility:  string(result.Visibility),
		Author: domain.NoteAuthor{
			ID:        result.Author.ID,
			Name:      result.Author.Name,
			AvatarURL: result.Author.AvatarURL,
		},
		UpdatedAt: result.UpdatedAt,
		CreatedAt: result.CreatedAt,
	})
}

// @Summary Delete a note by ID
// @Description Delete an existing note based on the provided ID
// @Tags note
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 "Successfully deleted a note by ID"
// @Router /notes/{id} [delete]
func (h *NoteHandler) Delete(ctx *fiber.Ctx) error {
	var req domain.NoteRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, ok := ctx.Locals("claims").(*domain.Claims)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "failed to retrieve claims from context")
	}
	if err := h.service.Delete(req.ID, *claims); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON("successfully deleted note by id")
}
