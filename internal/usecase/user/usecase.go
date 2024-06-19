package user

import (
	"context"
	"fmt"
	"tgstat/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
}

type Auth interface {
	HashPassword(password []byte) []byte
}

type UseCase struct {
	repo Repository
	auth Auth
}

func NewUseCase(repo Repository, auth Auth) *UseCase {
	return &UseCase{
		repo: repo,
		auth: auth,
	}
}

func (uc *UseCase) Create(ctx context.Context, user *entity.User) error {
	if user.Password == nil {
		return fmt.Errorf("invalid password")
	}

	hashedPassword := string(uc.auth.HashPassword([]byte(*user.Password)))
	user.Password = &hashedPassword

	if err := uc.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
