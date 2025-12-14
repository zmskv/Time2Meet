package handler

import (
	"net/http"
	"strconv"
	"time"

	"time2meet/internal/application/usecase/report"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	uc *report.UseCase
}

func NewReportHandler(uc *report.UseCase) *ReportHandler { return &ReportHandler{uc: uc} }

// @Summary Отчёт по продажам
// @Tags reports
// @Produce json
// @Param start query string true "Start date (YYYY-MM-DD)"
// @Param end query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} SalesReportRowSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/sales [get]
func (h *ReportHandler) Sales(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" || endStr == "" {
		RespondError(c, apperror.New(apperror.CodeValidation, "start and end are required (YYYY-MM-DD)", nil))
		return
	}
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid start date (YYYY-MM-DD)", err))
		return
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid end date (YYYY-MM-DD)", err))
		return
	}
	rows, err := h.uc.Sales(c.Request.Context(), start, end)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows)
}

// @Summary Статистика посещаемости по мероприятию
// @Tags reports
// @Produce json
// @Param event_id query string true "Event ID (UUID)"
// @Success 200 {array} AttendanceRowSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/attendance [get]
func (h *ReportHandler) Attendance(c *gin.Context) {
	eventIDStr := c.Query("event_id")
	eventID, err := valueobject.ParseUUID(eventIDStr)
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid event_id", err))
		return
	}
	rows, err := h.uc.Attendance(c.Request.Context(), eventID)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows)
}

// @Summary Популярные мероприятия
// @Tags analytics
// @Produce json
// @Param limit query int false "Limit"
// @Param days query int false "Days"
// @Success 200 {array} PopularEventRowSwagger
// @Failure 500 {object} ErrorResponse
// @Router /analytics/popular-events [get]
func (h *ReportHandler) Popular(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	rows, err := h.uc.Popular(c.Request.Context(), limit, days)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows)
}
