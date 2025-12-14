package handler

import (
	"net/http"
	"strconv"

	"time2meet/internal/application/usecase/user"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *user.UseCase
}

func NewUserHandler(uc *user.UseCase) *UserHandler { return &UserHandler{uc: uc} }

type CreateUserRequest struct {
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	Phone        string `json:"phone"`
	Role         string `json:"role" binding:"required"`
}

// @Summary Создать пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "Пользователь"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	id, err := h.uc.Create(c.Request.Context(), user.CreateUserInput(req))
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, IDResponse{ID: id.String()})
}

// @Summary Получить пользователя по id
// @Tags users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} UserSwagger
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	u, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

// @Summary Список пользователей
// @Tags users
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} UserSwagger
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	users, err := h.uc.List(c.Request.Context(), limit, offset)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

type UpdateUserRequest struct {
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	Phone        string `json:"phone"`
	Role         string `json:"role" binding:"required"`
	IsActive     bool   `json:"is_active"`
}

// @Summary Обновить пользователя
// @Tags users
// @Accept json
// @Param id path string true "User ID (UUID)"
// @Param body body UpdateUserRequest true "Поля пользователя"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := valueobject.ParseUUID(c.Param("id"))
	if err != nil {
		RespondError(c, apperror.New(apperror.CodeValidation, "invalid id", err))
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, err)
		return
	}
	err = h.uc.Update(c.Request.Context(), user.UpdateUserInput{
		ID:           id,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Role:         req.Role,
		IsActive:     req.IsActive,
	})
	if err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Удалить пользователя
// @Tags users
// @Param id path string true "User ID (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
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
