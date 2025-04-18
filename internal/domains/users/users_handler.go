package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/codepnw/gopher-social/internal/domains/commons"
	"github.com/codepnw/gopher-social/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateHandler(c *gin.Context)
	GetByIDHandler(c *gin.Context)
}

type handler struct {
	uc UserUsecase
}

func NewUserHandler(uc UserUsecase) UserHandler {
	return &handler{uc: uc}
}

func (h *handler) CreateHandler(c *gin.Context) {
	var payload UserReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	user, err := h.uc.Create(c, &payload)
	if err != nil {
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

	response.ResponseData(c, http.StatusCreated, user)
}

func (h *handler) GetByIDHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	user, err := h.uc.GetByID(c, id)
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			response.NotFoundResponse(c, err)
		default:
			log.Println("Here ??")
			response.InternalServerError(c, err)
		}
		return
	}

	response.ResponseData(c, http.StatusOK, user)
}