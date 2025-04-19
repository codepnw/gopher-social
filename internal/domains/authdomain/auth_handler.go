package authdomain

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/auth"
	"github.com/codepnw/gopher-social/internal/domains/commons"
	"github.com/codepnw/gopher-social/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type handler struct {
	uc     AuthUsecase
	config config.Config
	jwt    *auth.JWTAuthenticator
}

func NewAuthHandler(uc AuthUsecase, config config.Config, jwt *auth.JWTAuthenticator) AuthHandler {
	return &handler{
		uc:     uc,
		config: config,
		jwt:    jwt,
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

func (h *handler) Login(c *gin.Context) {
	var payload LoginUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	user, err := h.uc.GetUser(c, payload)
	if err != nil {
		switch err {
		case commons.ErrInvalidEmailPassword:
			response.BadRequestResponse(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(h.config.Auth.JWTExp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.JWTIss,
		"aud": h.config.Auth.JWTIss,
	}

	token, err := h.jwt.GenerateToken(claims)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.ResponseData(c, http.StatusOK, token)
}
