package userhandler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tgstat/internal/entity"
)

type UserUseCase interface {
	Create(ctx context.Context, user *entity.User) error
}

type Handler struct {
	user UserUseCase
}

func NewHandler(user UserUseCase) *Handler {
	return &Handler{
		user: user,
	}
}

func (h *Handler) Create(c *gin.Context) {
	req := new(entity.User)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := h.user.Create(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
}
