package handler

import (
	"net/http"
	"strconv"

	"time2meet/internal/application/usecase/venue"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	uc *venue.UseCase
}

func NewVenueHandler(uc *venue.UseCase) *VenueHandler { return &VenueHandler{uc: uc} }

type CreateVenueRequest struct {
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	City         string `json:"city" binding:"required"`
	Country      string `json:"country"`
	Capacity     int    `json:"capacity"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Website      string `json:"website"`
}

// @Summary Создать площадку
// @Tags venues
// @Accept json
// @Produce json
// @Param body body CreateVenueRequest true "Площадка"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /venues [post]
func (h *VenueHandler) CreateVenue(c *gin.Context) {
	var req CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	id, err := h.uc.CreateVenue(c.Request.Context(), venue.CreateVenueInput(req))
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, IDResponse{ID: id.String()})
}

// @Summary Получить площадку по id
// @Tags venues
// @Produce json
// @Param id path string true "Venue ID (UUID)"
// @Success 200 {object} VenueSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /venues/{id} [get]
func (h *VenueHandler) GetVenue(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	v, err := h.uc.GetVenue(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

// @Summary Список площадок
// @Tags venues
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} VenueSwagger
// @Failure 500 {object} ErrorResponse
// @Router /venues [get]
func (h *VenueHandler) ListVenues(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	venues, err := h.uc.ListVenues(c.Request.Context(), limit, offset)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, venues)
}

type UpdateVenueRequest struct {
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	City         string `json:"city" binding:"required"`
	Country      string `json:"country"`
	Capacity     int    `json:"capacity"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Website      string `json:"website"`
	IsActive     bool   `json:"is_active"`
}

// @Summary Обновить площадку
// @Tags venues
// @Accept json
// @Param id path string true "Venue ID (UUID)"
// @Param body body UpdateVenueRequest true "Поля площадки"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /venues/{id} [put]
func (h *VenueHandler) UpdateVenue(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	var req UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	err = h.uc.UpdateVenue(c.Request.Context(), venue.UpdateVenueInput{
		ID:           id,
		Name:         req.Name,
		Address:      req.Address,
		City:         req.City,
		Country:      req.Country,
		Capacity:     req.Capacity,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Website:      req.Website,
		IsActive:     req.IsActive,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Удалить площадку
// @Tags venues
// @Param id path string true "Venue ID (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /venues/{id} [delete]
func (h *VenueHandler) DeleteVenue(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	if err := h.uc.DeleteVenue(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

type CreateRoomRequest struct {
	Name        string         `json:"name" binding:"required"`
	Capacity    int            `json:"capacity"`
	Floor       *int           `json:"floor"`
	Equipment   map[string]any `json:"equipment"`
	HourlyRate  string         `json:"hourly_rate"`
	IsAvailable bool           `json:"is_available"`
}

// @Summary Создать помещение на площадке
// @Tags venues
// @Accept json
// @Produce json
// @Param id path string true "Venue ID (UUID)"
// @Param body body CreateRoomRequest true "Помещение"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /venues/{id}/rooms [post]
func (h *VenueHandler) CreateRoom(c *gin.Context) {
	venueID, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid venue id", err))
		return
	}
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	id, err := h.uc.CreateRoom(c.Request.Context(), venue.CreateRoomInput{
		VenueID:     venueID,
		Name:        req.Name,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		Equipment:   req.Equipment,
		HourlyRate:  req.HourlyRate,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, IDResponse{ID: id.String()})
}

// @Summary Список помещений площадки
// @Tags venues
// @Produce json
// @Param id path string true "Venue ID (UUID)"
// @Success 200 {array} RoomSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /venues/{id}/rooms [get]
func (h *VenueHandler) ListRooms(c *gin.Context) {
	venueID, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid venue id", err))
		return
	}
	rooms, err := h.uc.ListRooms(c.Request.Context(), venueID)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rooms)
}
