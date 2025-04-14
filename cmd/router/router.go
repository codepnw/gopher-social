package router

import (
	"github.com/codepnw/gopher-social/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (app *Application) Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	mid := middleware.InitMiddleware()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v := r.Group("/" + app.Config.App.ApiVersion)

	authorized := v.Group("/", gin.BasicAuth(gin.Accounts{
		app.Config.Auth.BasicUser: app.Config.Auth.BasicPassword,
	}))

	authorized.GET("/health", app.HealthCheckHandler)
	authorized.GET("/basic-logout", app.BasicAuthLogout)

	postsAuth := v.Group("/posts", mid.AuthTokenMiddleware())
	{
		postsAuth.POST("/", app.Store.Posts.CreatePostHandler)

		postsID := postsAuth.Group("/:id")
		{
			// Middleware
			// postsID.Use(app.Store.Posts.PostContextMiddleware())

			postsID.GET("/", app.Store.Posts.GetPostHandler)
			postsID.PATCH("/", app.Store.Posts.UpdatePostHandler)
			postsID.DELETE("/", app.Store.Posts.DeletePostHandler)
		}
	}

	users := v.Group("/users")
	{
		users.PUT("/activate/:token", app.Store.Users.ActivateUserHandler)

		usersAuth := users.Group("/:id", mid.AuthTokenMiddleware())
		{
			// Middleware
			// usersAuth.Use(app.Store.Users.UserContextMiddleware())

			usersAuth.GET("/", app.Store.Users.GetUserHandler)
			usersAuth.PUT("/follow", app.Store.Users.FollowUserHandler)
			usersAuth.PUT("/unfollow", app.Store.Users.UnfollowUserHandler)
		}
		usersAuth.GET("/feed", app.Store.Posts.GetUserFeedHandler)
	}

	auth := v.Group("/auth")
	{
		auth.POST("/register", app.Store.Users.RegisterUserHandler)
		auth.POST("/token", app.Store.Users.CreateTokenHandler)
	}

	return r
}
