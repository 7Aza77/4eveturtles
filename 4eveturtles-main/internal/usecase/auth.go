package usecase

import (
	"context"
	"errors"
	"goevent/internal/entity"
	"goevent/internal/repository"
	"goevent/pkg/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type Auth struct {
	repo         repository.AuthRepository
	tokenManager auth.TokenManager
	tokenTTL     time.Duration
}

func NewAuth(repo repository.AuthRepository, tokenManager auth.TokenManager, tokenTTL time.Duration) *Auth {
	return &Auth{
		repo:         repo,
		tokenManager: tokenManager,
		tokenTTL:     tokenTTL,
	}
}

func (u *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := entity.User{
		Email:    email,
		Password: string(passwordHash),
		Role:     entity.RoleStudent, // По умолчанию студент
	}

	return u.repo.CreateUser(ctx, user)
}

func (u *Auth) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found or invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return u.tokenManager.NewJWT(user.ID, string(user.Role), u.tokenTTL)
}
