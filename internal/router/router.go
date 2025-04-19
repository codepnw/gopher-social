package newrouter

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/auth"
	"github.com/codepnw/gopher-social/internal/domains/authdomain"
	"github.com/codepnw/gopher-social/internal/domains/feed"
	"github.com/codepnw/gopher-social/internal/domains/posts"
	"github.com/codepnw/gopher-social/internal/domains/users"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	DB     *sql.DB
	Config config.Config
	JWT    *auth.JWTAuthenticator
}

func (s *Routes) SetupRoutes() *gin.Engine {
	auth := authdomain.InitAuthDomain(s.DB, s.Config, s.JWT)
	post := posts.InitPostDomain(s.DB)
	user := users.InitUserDomain(s.DB, s.Config)
	feed := feed.InitFeedDomain(s.DB)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	version := s.Config.App.ApiVersion
	port := fmt.Sprintf(":%s", s.Config.App.Addr)

	// Basic Auth
	basicauth := r.Group("/", gin.BasicAuth(gin.Accounts{
		s.Config.Auth.BasicUser: s.Config.Auth.BasicPassword,
	}))
	{
		basicauth.GET(version+"/health", s.healthCheckHandler)
		basicauth.GET(version+"/logout", s.basicAuthLogout)
	}

	// Auth Routes
	authroutes := r.Group(version + "/auth")
	authroutes.POST("/register", auth.Register)
	authroutes.POST("/login", auth.Login)

	// Post Routes
	postroutes := r.Group(version + "/posts")
	postroutes.POST("/", post.CreatePostHandler)
	{
		postroutes.Use(post.PostContextMiddleware())
		postroutes.GET("/:id", post.GetPostHandler)
		postroutes.PATCH("/:id", post.UpdatePostHandler)
		postroutes.DELETE("/:id", post.DeletePostHandler)
	}

	// User Routes
	userroutes := r.Group(version + "/users")
	userroutes.POST("/", user.CreateHandler)
	userroutes.PUT("/activate/:token", user.ActivateHandler)
	{
		userroutes.Use(user.UserContextMiddleware())
		userroutes.GET("/:id", user.GetByIDHandler)
		userroutes.GET("/:id/follow", user.FollowUserHandler)
		userroutes.GET("/:id/unfollow", user.UnfollowUserHandler)
		userroutes.GET("/:id/feed", feed.GetUserFeedHandler)
	}

	r.Run(port)

	return r
}

func (s *Routes) healthCheckHandler(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)

	data := map[string]string{
		"status":  "ok",
		"env":     s.Config.App.Env,
		"version": s.Config.App.AppVersion,
		"user":    user,
	}

	c.JSON(http.StatusOK, data)
}

func (s *Routes) basicAuthLogout(c *gin.Context) {
	c.Header("WWW-Authenticate", `Basic realm="Please re-login"`)
	c.AbortWithStatus(http.StatusUnauthorized)
}
