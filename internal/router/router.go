package newrouter

import (
	"database/sql"
	"fmt"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/domains/auth"
	"github.com/codepnw/gopher-social/internal/domains/feed"
	"github.com/codepnw/gopher-social/internal/domains/posts"
	"github.com/codepnw/gopher-social/internal/domains/users"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	DB     *sql.DB
	Config config.Config
}

func (s *Routes) SetupRoutes() *gin.Engine {
	auth := auth.InitAuthDomain(s.DB, s.Config)
	post := posts.InitPostDomain(s.DB)
	user := users.InitUserDomain(s.DB, s.Config)
	feed := feed.InitFeedDomain(s.DB)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	version := s.Config.App.ApiVersion
	port := fmt.Sprintf(":%s", s.Config.App.Addr)

	// Auth Routes
	authroutes := r.Group(version + "/auth")

	authroutes.POST("/register", auth.Register)

	// Post Routes
	postroutes := r.Group(version + "/posts")

	postroutes.POST("/", post.CreatePostHandler)
	postroutes.Use(post.PostContextMiddleware())
	postroutes.GET("/:id", post.GetPostHandler)
	postroutes.PATCH("/:id", post.UpdatePostHandler)
	postroutes.DELETE("/:id", post.DeletePostHandler)

	// User Routes
	userroutes := r.Group(version + "/users")

	userroutes.POST("/", user.CreateHandler)
	userroutes.GET("/:id", user.GetByIDHandler)

	userroutes.GET("/:id/feed", feed.GetUserFeedHandler)

	r.Run(port)

	return r
}
