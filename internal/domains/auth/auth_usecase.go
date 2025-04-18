package auth

import (
	"context"
	"time"

	"github.com/codepnw/gopher-social/internal/domains/users"
)

type AuthUsecase interface {
	Register(ctx context.Context, payload *RegisterUserPayload, token string, exp time.Duration) error
}

type usecase struct {
	userRepo users.UserUsecase
}

func NewAuthUsecase(userRepo users.UserUsecase) AuthUsecase {
	return &usecase{userRepo: userRepo}
}

func (uc *usecase) Register(ctx context.Context, payload *RegisterUserPayload, token string, exp time.Duration) error {
	user := &users.UserReq{
		Email:    payload.Email,
		Username: payload.Username,
	}

	if err := user.HashPassword(payload.Password); err != nil {
		return err
	}

	if err := uc.userRepo.CreateAndInvite(ctx, user, token, exp); err != nil {
		return err
	}

	return nil
}
