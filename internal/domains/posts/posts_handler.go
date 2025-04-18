package posts

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/codepnw/gopher-social/internal/domains/commons"
	"github.com/codepnw/gopher-social/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type PostHandler interface {
	CreatePostHandler(c *gin.Context)
	GetPostHandler(c *gin.Context)
	UpdatePostHandler(c *gin.Context)
	DeletePostHandler(c *gin.Context)
	PostContextMiddleware() gin.HandlerFunc
}

type handler struct {
	uc PostUsecase
}

func NewPostHandler(uc PostUsecase) PostHandler {
	return &handler{uc: uc}
}

func (h *handler) CreatePostHandler(c *gin.Context) {
	var payload CreatePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	post, err := h.uc.Create(c, &payload)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.ResponseData(c, http.StatusOK, post)
}

func (h *handler) GetPostHandler(c *gin.Context) {
	post := h.getPostContext(c)

	// TODO: get comments later

	response.ResponseData(c, http.StatusOK, post)
}

func (h *handler) UpdatePostHandler(c *gin.Context) {
	post := h.getPostContext(c)

	var payload UpdatePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	p, err := h.uc.Update(c, post.ID, &payload)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.ResponseData(c, http.StatusOK, p)
}

func (h *handler) DeletePostHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.uc.Delete(c, id); err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.ResponseData(c, http.StatusNoContent, nil)
}

func (h *handler) PostContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		post, err := h.uc.GetByID(c, id)
		if err != nil {
			switch {
			case errors.Is(err, commons.ErrNotFound):
				response.NotFoundResponse(c, err)
			default:
				response.InternalServerError(c, err)
			}
			return
		}

		c.Set(commons.ContextPostKey, post)
		c.Next()
	}
}

func (h *handler) getPostContext(c *gin.Context) *Post {
	post, _ := c.Get(commons.ContextPostKey)
	return post.(*Post)
}
