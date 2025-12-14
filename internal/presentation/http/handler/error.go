package handler

import (
	"errors"
	"net/http"

	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

func RespondError(c *gin.Context, err error) {
	var ae *apperror.AppError
	if errors.As(err, &ae) {
		c.JSON(mapCodeToStatus(ae.Code), gin.H{
			"code":    ae.Code,
			"message": ae.Message,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    apperror.CodeInternal,
		"message": "internal error",
	})
}

func mapCodeToStatus(code apperror.Code) int {
	switch code {
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict:
		return http.StatusConflict
	case apperror.CodeValidation:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	case apperror.CodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}


