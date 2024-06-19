package channel

import (
	"context"
	"fmt"
	"strings"
	"tgstat/internal/entity"
	"time"
)

type ChannelRepository interface {
	GetAll(ctx context.Context) ([]*entity.Channel, error)
	GetByID(ctx context.Context, id int) (*entity.Channel, error)
	Create(ctx context.Context, channel *entity.Channel) error
	UpdateUpdatedTime(ctx context.Context, channel *entity.Channel) error
}

type ChannelService interface {
	GetInfo(ctx context.Context, channelID string) (*entity.Channel, error)
}

type PostService interface {
	SyncAllPosts(ctx context.Context, channelID string) error
}

type UseCase struct {
	chanRepo    ChannelRepository
	chanService ChannelService
	post        PostService
}

func NewUseCase(chanRepo ChannelRepository, chanService ChannelService, post PostService) *UseCase {
	return &UseCase{
		chanRepo:    chanRepo,
		chanService: chanService,
		post:        post,
	}
}

func (uc *UseCase) GetAllChannels(ctx context.Context) ([]*entity.Channel, error) {
	return uc.chanRepo.GetAll(ctx)
}

func (uc *UseCase) GetChannelByID(ctx context.Context, id int) (*entity.Channel, error) {
	return uc.chanRepo.GetByID(ctx, id)
}

func (uc *UseCase) CreateChannel(ctx context.Context, channelID string) (*entity.Channel, error) {
	channel, err := uc.chanService.GetInfo(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(channel.Image100, "//") {
		channel.Image100 = channel.Image100[2:]
	}
	if strings.HasPrefix(channel.Image640, "//") {
		channel.Image640 = channel.Image640[2:]
	}

	if err := uc.chanRepo.Create(ctx, channel); err != nil {
		return nil, err
	}

	go func() {
		ctx = context.Background()
		if err := uc.post.SyncAllPosts(ctx, channel.Username); err != nil {
			fmt.Println(err)
		}
		channel.LastUpdated = time.Now().Unix()
		if err := uc.chanRepo.UpdateUpdatedTime(ctx, channel); err != nil {
			fmt.Println(err)
		}
	}()

	return channel, nil
}
