package authdomain

import (
	"context"
	"log"
	"time"

	"github.com/codepnw/gopher-social/internal/domains/commons"
	"github.com/codepnw/gopher-social/internal/domains/users"
)

type AuthUsecase interface {
	Register(ctx context.Context, payload *RegisterUserPayload, token string, exp time.Duration) error
	GetUser(ctx context.Context, req LoginUserPayload) (*users.User, error)
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

func (uc *usecase) GetUser(ctx context.Context, req LoginUserPayload) (*users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
	log.Println("GetUSER=========1", err)

		return nil, commons.ErrInvalidEmailPassword
	}

	if err = user.ComparePassword(req.Password); err != nil {
	log.Println("GetUSER=========2", user)

		return nil, commons.ErrInvalidEmailPassword
	}

	log.Println("GetUSER=========", user)

	return user, nil
}
