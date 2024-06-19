package authorization

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"tgstat/internal/entity"
	"time"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type UseCase struct {
	userRepo   UserRepository
	secretKey  []byte
	expireTime time.Duration
}

func NewUseCase(userRepo UserRepository, secretKey string, expireTime time.Duration) *UseCase {
	return &UseCase{
		userRepo:   userRepo,
		secretKey:  []byte(secretKey),
		expireTime: expireTime,
	}
}

func (uc *UseCase) SignIn(ctx context.Context, username, password string) (*entity.User, string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if user.Password == nil {
		return nil, "", fmt.Errorf("user has no password")
	}

	if !uc.ComparePasswords([]byte(*user.Password), []byte(password)) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	user.Password = (*string)(nil)

	token, err := uc.GenerateJWTToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (uc *UseCase) GenerateJWTToken(id uuid.UUID, username string) (string, error) {
	claims := &entity.Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * 8),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(uc.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uc *UseCase) VerifyToken(tokenString string) (*entity.User, error) {
	claims := new(entity.Claims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}

		return uc.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	user := &entity.User{
		ID:       claims.ID,
		Username: claims.Username,
	}

	return user, nil
}

func (uc *UseCase) HashPassword(password []byte) []byte {
	pass, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return pass
}

func (uc *UseCase) ComparePasswords(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}
