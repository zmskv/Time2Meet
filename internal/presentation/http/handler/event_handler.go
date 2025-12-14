package handler

import (
	"net/http"
	"strconv"

	"time2meet/internal/application/usecase/event"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	uc *event.UseCase
}

func NewEventHandler(uc *event.UseCase) *EventHandler { return &EventHandler{uc: uc} }

type CreateEventRequest struct {
	OrganizerID     string `json:"organizer_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	Status          string `json:"status" binding:"required"`
	IsPublic        bool   `json:"is_public"`
	MaxParticipants *int   `json:"max_participants"`
	CoverImage      string `json:"cover_image"`
}

// @Summary Создать мероприятие
// @Tags events
// @Accept json
// @Produce json
// @Param body body CreateEventRequest true "Мероприятие"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events [post]
func (h *EventHandler) Create(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	orgID, err := valueobject.ParseUUID(req.OrganizerID)
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid organizer_id", err))
		return
	}
	id, err := h.uc.Create(c.Request.Context(), event.CreateEventInput{
		OrganizerID:     orgID,
		Title:           req.Title,
		Description:     req.Description,
		Status:          req.Status,
		IsPublic:        req.IsPublic,
		MaxParticipants: req.MaxParticipants,
		CoverImage:      req.CoverImage,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, IDResponse{ID: id.String()})
}

// @Summary Получить мероприятие по id
// @Tags events
// @Produce json
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} EventSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events/{id} [get]
func (h *EventHandler) Get(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	e, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, e)
}

// @Summary Список мероприятий
// @Tags events
// @Produce json
// @Param organizer_id query string false "Organizer ID (UUID)"
// @Param status query string false "Status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} EventSwagger
// @Failure 500 {object} ErrorResponse
// @Router /events [get]
func (h *EventHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var organizerID *valueobject.UUID
	if v := c.Query("organizer_id"); v != "" {
		id, err := valueobject.ParseUUID(v)
		if err == nil {
			organizerID = &id
		}
	}
	var status *string
	if v := c.Query("status"); v != "" {
		status = &v
	}
	events, err := h.uc.List(c.Request.Context(), organizerID, status, limit, offset)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, events)
}

type UpdateEventRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	Status          string `json:"status" binding:"required"`
	IsPublic        bool   `json:"is_public"`
	MaxParticipants *int   `json:"max_participants"`
	CoverImage      string `json:"cover_image"`
}

// @Summary Обновить мероприятие
// @Tags events
// @Accept json
// @Param id path string true "Event ID (UUID)"
// @Param body body UpdateEventRequest true "Поля мероприятия"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events/{id} [put]
func (h *EventHandler) Update(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	err = h.uc.Update(c.Request.Context(), event.UpdateEventInput{
		ID:              id,
		Title:           req.Title,
		Description:     req.Description,
		Status:          req.Status,
		IsPublic:        req.IsPublic,
		MaxParticipants: req.MaxParticipants,
		CoverImage:      req.CoverImage,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Удалить мероприятие
// @Tags events
// @Param id path string true "Event ID (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events/{id} [delete]
func (h *EventHandler) Delete(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Отменить мероприятие
// @Tags events
// @Param id path string true "Event ID (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events/{id}/cancel [post]
func (h *EventHandler) Cancel(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	if err := h.uc.Cancel(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
