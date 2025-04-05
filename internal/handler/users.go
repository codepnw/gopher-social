package handler

import (
	"net/http"
	"strconv"

	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/codepnw/gopher-social/internal/repository"
	"github.com/gin-gonic/gin"
)

const userCtxKey string = "user"

type UserHandler interface {
	CreateUserHandler(c *gin.Context)
	GetUserHandler(c *gin.Context)
	FollowUserHandler(c *gin.Context)
	UnfollowUserHandler(c *gin.Context)
	UserContextMiddleware() gin.HandlerFunc
}

type userHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) UserHandler {
	return &userHandler{repo: repo}
}

func (h *userHandler) CreateUserHandler(c *gin.Context) {
	var req entity.UserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestResponse(c, err)
		return
	}

	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.repo.Create(c, &user); err != nil {
		internalServerError(c, err)
		return
	}

	responseData(c, http.StatusCreated, user)
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	user := getUserFromContext(c)

	responseData(c, http.StatusOK, user)
}

func (h *userHandler) FollowUserHandler(c *gin.Context) {
	followerUser := getUserFromContext(c)

	// TODO: revert back to auth userID
	var payload entity.FollowUser
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequestResponse(c, err)
		return
	}

	if err := h.repo.Follow(c, followerUser.ID, payload.UserID); err != nil {
		switch err {
		case repository.ErrConflict:
			conflictResponse(c, err)
			return
		default:
			internalServerError(c, err)
			return
		}
	}

	responseData(c, http.StatusNoContent, nil)
}

func (h *userHandler) UnfollowUserHandler(c *gin.Context) {
	unfollowedUser := getUserFromContext(c)

	// TODO: revert back to auth userID
	var payload entity.FollowUser
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequestResponse(c, err)
		return
	}

	if err := h.repo.Unfollow(c, unfollowedUser.ID, payload.UserID); err != nil {
		internalServerError(c, err)
		return
	}

	responseData(c, http.StatusNoContent, nil)
}

func (h *userHandler) UserContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			badRequestResponse(c, err)
			return
		}

		user, err := h.repo.GetByID(c, id)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				notFoundResponse(c, err)
				return
			default:
				internalServerError(c, err)
				return
			}
		}

		c.Set(userCtxKey, user)
		c.Next()
	}
}

func getUserFromContext(c *gin.Context) *entity.User {
	user, _ := c.Get(userCtxKey)
	return user.(*entity.User)
}
