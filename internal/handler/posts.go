package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/codepnw/gopher-social/internal/repository"
	"github.com/codepnw/gopher-social/internal/utils"
	"github.com/gin-gonic/gin"
)

const postCtxKey string = "post"

type PostsHandler interface {
	CreatePostHandler(c *gin.Context)
	GetPostHandler(c *gin.Context)
	UpdatePostHandler(c *gin.Context)
	DeletePostHandler(c *gin.Context)
	PostContextMiddleware() gin.HandlerFunc
}

type postHandler struct {
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
}

func NewPostsHandler(postRepo repository.PostRepository, commentRepo repository.CommentRepository) PostsHandler {
	return &postHandler{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (h *postHandler) CreatePostHandler(c *gin.Context) {
	var payload entity.CreatePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequestResponse(c, err)
		return
	}

	post := &entity.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Change after auth
		UserID: 1,
	}

	if err := h.postRepo.Create(c, post); err != nil {
		internalServerError(c, err)
		return
	}

	responseData(c, http.StatusCreated, post)
}

func (h *postHandler) GetPostHandler(c *gin.Context) {
	post := getPostFromContext(c)

	// Comments
	comments, err := h.commentRepo.GetByPostID(c, post.ID)
	if err != nil {
		internalServerError(c, err)
		return
	}

	post.Comments = comments

	responseData(c, http.StatusOK, post)
}

func (h *postHandler) DeletePostHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	if err := h.postRepo.Delete(c, idInt); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			notFoundResponse(c, err)
		default:
			internalServerError(c, err)
		}
		return
	}

	responseData(c, http.StatusNoContent, "post deleted")
}

func (h *postHandler) UpdatePostHandler(c *gin.Context) {
	post := getPostFromContext(c)

	var payload entity.UpdatePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequestResponse(c, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	post.UpdatedAt = utils.TimeString()

	if err := h.postRepo.Update(c, post); err != nil {
		internalServerError(c, err)
		return
	}

	responseData(c, http.StatusOK, post)
}

func (h *postHandler) PostContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			internalServerError(c, err)
			return
		}

		post, err := h.postRepo.GetByID(c, idInt)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				notFoundResponse(c, err)
			default:
				internalServerError(c, err)
			}
			return
		}

		// ctx := context.WithValue(c, "post", post)
		c.Set(postCtxKey, post)
		c.Next()
	}
}

func getPostFromContext(c *gin.Context) *entity.Post {
	post, _ := c.Get(postCtxKey)
	return post.(*entity.Post)
}
