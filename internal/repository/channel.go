package repository

import (
	"context"
	"github.com/uptrace/bun"
	"tgstat/internal/entity"
)

type ChannelRepository struct {
	db *bun.DB
}

func NewChannelRepository(db *bun.DB) *ChannelRepository {
	return &ChannelRepository{
		db: db,
	}
}

func (r *ChannelRepository) GetAll(ctx context.Context) ([]*entity.Channel, error) {
	channels := make([]*entity.Channel, 0)
	if err := r.db.NewSelect().Model(&channels).Scan(ctx); err != nil {
		return nil, err
	}
	return channels, nil
}

func (r *ChannelRepository) GetByID(ctx context.Context, id int) (*entity.Channel, error) {
	channel := new(entity.Channel)
	if err := r.db.NewSelect().Model(channel).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return channel, nil
}

func (r *ChannelRepository) Create(ctx context.Context, channel *entity.Channel) error {
	_, err := r.db.NewInsert().Model(channel).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChannelRepository) UpdateUpdatedTime(ctx context.Context, channel *entity.Channel) error {
	_, err := r.db.NewUpdate().Model(channel).
		Set("last_updated = ?", channel.LastUpdated).Where("id = ?", channel.Id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
