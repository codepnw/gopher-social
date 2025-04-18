package feed

import (
	"net/http"
	"strconv"

	"github.com/codepnw/gopher-social/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type FeedHandler interface {
	GetUserFeedHandler(c *gin.Context)
}

type handler struct {
	uc FeedUsecase
}

func NewFeedHandler(uc FeedUsecase) FeedHandler {
	return &handler{uc: uc}
}

func (h *handler) GetUserFeedHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	var fq PaginatedFeedQuery

	fq, err = fq.Parse(c)
	if err != nil {
		response.BadRequestResponse(c, err)
		return
	}

	feed, err := h.uc.GetUserFeed(c, userID, fq)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	// TODO: feed empty fix later

	response.ResponseData(c, http.StatusOK, feed)
}
