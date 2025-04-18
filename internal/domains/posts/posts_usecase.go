package posts

import (
	"context"
	"database/sql"
	"time"

	"github.com/codepnw/gopher-social/internal/domains/commons"
)

type PostUsecase interface {
	Create(ctx context.Context, post *CreatePostPayload) (*Post, error)
	GetByID(ctx context.Context, postID int64) (*Post, error)
	Update(ctx context.Context, id int64, newPost *UpdatePostPayload) (*Post, error)
	Delete(ctx context.Context, postID int64) error
}

type usecase struct {
	repo PostRepository
}

func NewPostUsecase(repo PostRepository) PostUsecase {
	return &usecase{repo: repo}
}

func (uc *usecase) Create(ctx context.Context, post *CreatePostPayload) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	p := &Post{
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
		// TODO: change after
		UserID: 58,
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return &Post{}, err
	}

	return p, nil
}

func (uc *usecase) Update(ctx context.Context, id int64, newPost *UpdatePostPayload) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	var post Post

	if newPost.Title != nil {
		post.Title = *newPost.Title
	}

	if newPost.Content != nil {
		post.Content = *newPost.Content
	}

	post.ID = id
	post.UpdatedAt = time.Now().String()

	if err := uc.repo.Update(ctx, &post); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &Post{}, commons.ErrNotFound
		default:
			return &Post{}, err
		}
	}

	return &post, nil
}

func (uc *usecase) GetByID(ctx context.Context, postID int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	post, err := uc.repo.GetByID(ctx, postID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, commons.ErrNotFound
		default:
			return nil, err
		}
	}

	return post, err
}

func (uc *usecase) Delete(ctx context.Context, postID int64) error {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	if err := uc.repo.Delete(ctx, postID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return commons.ErrNotFound
		default:
			return err
		}
	}

	return nil
}
