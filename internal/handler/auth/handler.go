package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tgstat/internal/entity"
)

type Authorization interface {
	SignIn(ctx context.Context, username, password string) (*entity.User, string, error)
}

type Handler struct {
	auth Authorization
}

func NewHandler(auth Authorization) *Handler {
	return &Handler{
		auth: auth,
	}
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	User  *entity.User `json:"user"`
	Token string       `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var req signInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user, token, err := h.auth.SignIn(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res := &signInResponse{
		User:  user,
		Token: token,
	}

	c.JSON(http.StatusOK, res)
}
