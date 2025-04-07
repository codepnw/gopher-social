package entity

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type CreatePostPayload struct {
	Title   string   `json:"title" binding:"required,max=100"`
	Content string   `json:"content" binding:"required,max=300"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" binding:"omitempty,max=100"`
	Content *string `json:"content" binding:"omitempty,max=300"`
}

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" binding:"gte=1,lte=20"`
	Offset int      `json:"offset" binding:"gte=0"`
	Sort   string   `json:"sort" binding:"oneof=asc desc"`
	Tags   []string `json:"tags" binding:"max=5"`
	Search string   `json:"search" binding:"max=100"` // title, content
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(c *gin.Context) (PaginatedFeedQuery, error) {
	limit := c.Query("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := c.Query("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = o
	}

	sort := c.Query("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := c.Query("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := c.Query("search")
	if search != "" {
		fq.Search = search
	}

	since := c.Query("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := c.Query("until")
	if until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(since string) string {
	t, err := time.Parse(time.DateTime, since)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}