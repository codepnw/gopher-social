package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/codepnw/gopher-social/internal/repository"
	"github.com/codepnw/gopher-social/internal/utils/logger"
	"github.com/codepnw/gopher-social/internal/utils/mailer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const userCtxKey string = "user"

type UserHandler interface {
	RegisterUserHandler(c *gin.Context)
	ActivateUserHandler(c *gin.Context)
	GetUserHandler(c *gin.Context)
	FollowUserHandler(c *gin.Context)
	UnfollowUserHandler(c *gin.Context)
	UserContextMiddleware() gin.HandlerFunc
}

type userHandler struct {
	cfg    config.Config
	repo   repository.UserRepository
	mailer mailer.MailtrapClient
}

func NewUserHandler(cfg config.Config, repo repository.UserRepository, mailer mailer.MailtrapClient) UserHandler {
	return &userHandler{
		cfg:    cfg,
		repo:   repo,
		mailer: mailer,
	}
}

func (h *userHandler) RegisterUserHandler(c *gin.Context) {
	var payload entity.RegisterUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequestResponse(c, err)
		return
	}

	user := &entity.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.HashPassword(payload.Password); err != nil {
		internalServerError(c, err)
		return
	}

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// create invite user
	err := h.repo.CreateAndInvite(c, user, hashToken, h.cfg.Mail.Exp)
	if err != nil {
		switch err {
		case repository.ErrDuplicateEmail:
			badRequestResponse(c, err)
		case repository.ErrDuplicateUsername:
			badRequestResponse(c, err)
		default:
			internalServerError(c, err)
		}
		return
	}

	userToken := entity.UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", h.cfg.App.FrontendURL, plainToken)

	isProdEnv := h.cfg.App.Env == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// send mail
	_, err = h.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		logger.Error(c, "error sending welcome email", err)

		// rollback user creation if email fails
		if err := h.repo.Delete(c, user.ID); err != nil {
			logger.Error(c, "error deleting user", err)
		}

		internalServerError(c, err)
		return
	}

	responseData(c, http.StatusCreated, userToken)
}

func (h *userHandler) ActivateUserHandler(c *gin.Context) {
	token := c.Param("token")

	if err := h.repo.Activate(c, token); err != nil {
		switch err {
		case repository.ErrNotFound:
			badRequestResponse(c, err)
		default:
			internalServerError(c, err)
		}
		return
	}

	responseData(c, http.StatusNoContent, nil)
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
