package feed

import (
	"context"

	"github.com/codepnw/gopher-social/internal/domains/commons"
)

type FeedUsecase interface {
	GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error)
}

type usecase struct {
	repo FeedRepository
}

func NewFeedUsecase(repo FeedRepository) FeedUsecase {
	return &usecase{repo: repo}
}

func (uc *usecase) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	return uc.repo.GetUserFeed(ctx, userID, fq)
}
