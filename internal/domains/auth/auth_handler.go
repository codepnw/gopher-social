package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/domains/commons"
	"github.com/codepnw/gopher-social/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Register(c *gin.Context)
	// Login(c *gin.Context)
}

type handler struct {
	uc     AuthUsecase
	config config.Config
}

func NewAuthHandler(uc AuthUsecase, config config.Config) AuthHandler {
	return &handler{
		uc:     uc,
		config: config,
	}
}

func (h *handler) Register(c *gin.Context) {
	var payload RegisterUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := h.uc.Register(c, &payload, hashToken, h.config.Auth.JWTExp); err != nil {
		switch err {
		case commons.ErrDuplicateEmail:
			response.BadRequestResponse(c, err)
		case commons.ErrDuplicateUsername:
			response.BadRequestResponse(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	// TODO: Mail, Generate Token

	response.ResponseData(c, http.StatusCreated, gin.H{"token": plainToken})
}
