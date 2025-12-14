package handler

import (
	"net/http"

	"time2meet/internal/application/port/batchimport"
	"time2meet/internal/application/usecase/batch"
	"time2meet/internal/domain/valueobject"
	"time2meet/internal/presentation/http/middleware"
	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type BatchHandler struct {
	uc *batch.UseCase
}

func NewBatchHandler(uc *batch.UseCase) *BatchHandler { return &BatchHandler{uc: uc} }

type importUsersRequest struct {
	ContinueOnError bool                       `json:"continue_on_error"`
	Items           []batchimport.ImportUsersItem `json:"items" binding:"required"`
}

// @Summary Батч-импорт пользователей
// @Tags batch
// @Accept json
// @Produce json
// @Param X-User-Id header string false "User ID (UUID) for audit"
// @Param body body importUsersRequest true "Пакет пользователей"
// @Success 200 {object} batchimport.Result
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /batch/import/users [post]
func (h *BatchHandler) ImportUsers(c *gin.Context) {
	var req importUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid body", err))
		return
	}
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	userID, _ := uidAny.(valueobject.UUID)
	ipAny, _ := c.Get(middleware.CtxIPKey)
	ip, _ := ipAny.(string)

	out, err := h.uc.ImportUsers(c.Request.Context(), batch.ImportUsersInput{
		UserID:          userID,
		IP:              ip,
		ContinueOnError: req.ContinueOnError,
		Items:           req.Items,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

type importEventsRequest struct {
	ContinueOnError bool                        `json:"continue_on_error"`
	Items           []batchimport.ImportEventsItem `json:"items" binding:"required"`
}

// @Summary Батч-импорт мероприятий
// @Tags batch
// @Accept json
// @Produce json
// @Param X-User-Id header string false "User ID (UUID) for audit"
// @Param body body importEventsRequest true "Пакет мероприятий"
// @Success 200 {object} batchimport.Result
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /batch/import/events [post]
func (h *BatchHandler) ImportEvents(c *gin.Context) {
	var req importEventsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid body", err))
		return
	}
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	userID, _ := uidAny.(valueobject.UUID)
	ipAny, _ := c.Get(middleware.CtxIPKey)
	ip, _ := ipAny.(string)

	out, err := h.uc.ImportEvents(c.Request.Context(), batch.ImportEventsInput{
		UserID:          userID,
		IP:              ip,
		ContinueOnError: req.ContinueOnError,
		Items:           req.Items,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

type importTicketsRequest struct {
	ContinueOnError bool                         `json:"continue_on_error"`
	Items           []batchimport.ImportTicketsItem `json:"items" binding:"required"`
}

// @Summary Батч-импорт билетов
// @Tags batch
// @Accept json
// @Produce json
// @Param X-User-Id header string false "User ID (UUID) for audit"
// @Param body body importTicketsRequest true "Пакет билетов"
// @Success 200 {object} batchimport.Result
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /batch/import/tickets [post]
func (h *BatchHandler) ImportTickets(c *gin.Context) {
	var req importTicketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid body", err))
		return
	}
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	userID, _ := uidAny.(valueobject.UUID)
	ipAny, _ := c.Get(middleware.CtxIPKey)
	ip, _ := ipAny.(string)

	out, err := h.uc.ImportTickets(c.Request.Context(), batch.ImportTicketsInput{
		UserID:          userID,
		IP:              ip,
		ContinueOnError: req.ContinueOnError,
		Items:           req.Items,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
