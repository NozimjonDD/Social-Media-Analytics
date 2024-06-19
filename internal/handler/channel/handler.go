package channelhandler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tgstat/internal/entity"
)

type ChannelUseCase interface {
	GetAllChannels(ctx context.Context) ([]*entity.Channel, error)
	GetChannelByID(ctx context.Context, id int) (*entity.Channel, error)
	CreateChannel(ctx context.Context, channelID string) (*entity.Channel, error)
}

type PostUseCase interface {
	SyncAllPosts(ctx context.Context, channelID string) error
}

type Handler struct {
	channel ChannelUseCase
	post    PostUseCase
}

func NewHandler(channel ChannelUseCase, post PostUseCase) *Handler {
	return &Handler{
		channel: channel,
		post:    post,
	}
}

func (h *Handler) GetAllChannels(c *gin.Context) {
	channels, err := h.channel.GetAllChannels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, channels)
}

func (h *Handler) GetChannelByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}

	channel, err := h.channel.GetChannelByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, channel)
}

type channelCreateRequest struct {
	ChannelID string `json:"channel_id"`
}

func (h *Handler) CreateChannel(c *gin.Context) {
	req := new(channelCreateRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	channel, err := h.channel.CreateChannel(c.Request.Context(), req.ChannelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	//err = h.post.SyncAllPosts(c.Request.Context(), req.ChannelID)
	//if err != nil {
	//	if err.Error() == "wrong param offset: max count results - 1000" {
	//		c.JSON(http.StatusOK, channel)
	//	}
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	c.Abort()
	//	return
	//}

	c.JSON(http.StatusOK, channel)
}
