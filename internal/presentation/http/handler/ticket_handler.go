package handler

import (
	"net/http"
	"strconv"

	"time2meet/internal/application/usecase/ticket"
	"time2meet/internal/domain/valueobject"
	"time2meet/internal/presentation/http/middleware"
	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	purchase *ticket.PurchaseUseCase
	tickets  *ticket.TicketUseCase
	validate *ticket.ValidateUseCase
}

func NewTicketHandler(purchase *ticket.PurchaseUseCase, tickets *ticket.TicketUseCase, validate *ticket.ValidateUseCase) *TicketHandler {
	return &TicketHandler{purchase: purchase, tickets: tickets, validate: validate}
}

type PurchaseTicketRequest struct {
	TicketTypeID string `json:"ticket_type_id" binding:"required"`
	QRCode       string `json:"qr_code" binding:"required"`
	AmountPaid   string `json:"amount_paid" binding:"required"`
	Currency     string `json:"currency"`
}

// @Summary Купить билет (транзакция)
// @Tags tickets
// @Accept json
// @Produce json
// @Param X-User-Id header string false "User ID (UUID) for audit"
// @Param body body PurchaseTicketRequest true "Покупка"
// @Success 201 {object} TicketIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tickets/purchase [post]
func (h *TicketHandler) Purchase(c *gin.Context) {
	var req PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}

	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	userID, _ := uidAny.(valueobject.UUID)
	ipAny, _ := c.Get(middleware.CtxIPKey)
	ip, _ := ipAny.(string)

	ttid, err := valueobject.ParseUUID(req.TicketTypeID)
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid ticket_type_id", err))
		return
	}

	out, err := h.purchase.Purchase(c.Request.Context(), ticket.PurchaseInput{
		UserID:       userID,
		IP:           ip,
		TicketTypeID: ttid,
		QRCode:       req.QRCode,
		AmountPaid:   req.AmountPaid,
		Currency:     req.Currency,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, TicketIDResponse{TicketID: out.TicketID.String()})
}

// @Summary Получить билет по id
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID (UUID)"
// @Success 200 {object} TicketSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tickets/{id} [get]
func (h *TicketHandler) Get(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	t, err := h.tickets.Get(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, t)
}

// @Summary Список билетов по покупателю
// @Tags tickets
// @Produce json
// @Param buyer_id query string true "Buyer ID (UUID)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} TicketSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tickets [get]
func (h *TicketHandler) ListByBuyer(c *gin.Context) {
	buyerIDStr := c.Query("buyer_id")
	buyerID, err := valueobject.ParseUUID(buyerIDStr)
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid buyer_id", err))
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	rows, err := h.tickets.ListByBuyer(c.Request.Context(), buyerID, limit, offset)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows)
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// @Summary Обновить статус билета
// @Tags tickets
// @Accept json
// @Param id path string true "Ticket ID (UUID)"
// @Param body body UpdateTicketStatusRequest true "Статус"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tickets/{id}/status [patch]
func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	if err := h.tickets.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Удалить билет
// @Tags tickets
// @Param id path string true "Ticket ID (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tickets/{id} [delete]
func (h *TicketHandler) Delete(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	if err := h.tickets.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Валидация (использование) билета
// @Tags tickets
// @Param X-User-Id header string true "User ID (UUID) for audit"
// @Param id path string true "Ticket ID (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tickets/{id}/validate [post]
func (h *TicketHandler) Validate(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	userID, _ := uidAny.(valueobject.UUID)
	ipAny, _ := c.Get(middleware.CtxIPKey)
	ip, _ := ipAny.(string)

	if err := h.validate.Validate(c.Request.Context(), ticket.ValidateInput{
		UserID:   userID,
		IP:       ip,
		TicketID: id,
	}); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
