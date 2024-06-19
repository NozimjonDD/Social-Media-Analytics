package post

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tgstat/internal/entity"
)

type PostUseCase interface {
	GetAllPosts(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error)
	GetAllForWeek(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error)
	GetFrequencyWords(ctx context.Context, channelID int) ([]*entity.FrequencyWords, error)
}

type Handler struct {
	channel PostUseCase
}

func NewHandler(channel PostUseCase) *Handler {
	return &Handler{
		channel: channel,
	}
}

func (h *Handler) GetAllPost(c *gin.Context) {
	var filter entity.PostFilter
	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		c.Abort()
		return
	}
	posts, err := h.channel.GetAllPosts(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetAllForWeek(c *gin.Context) {
	var filter entity.PostFilter
	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		c.Abort()
		return
	}
	posts, err := h.channel.GetAllForWeek(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetFrequencyWords(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		c.Abort()
		return
	}
	posts, err := h.channel.GetFrequencyWords(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, posts)
}
