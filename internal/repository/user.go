package repository

import (
	"context"
	"github.com/uptrace/bun"
	"tgstat/internal/entity"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	err := r.db.NewSelect().Model(user).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
