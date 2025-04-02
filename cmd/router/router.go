package router

import "github.com/gin-gonic/gin"

func (app *Application) Routes() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v := r.Group("/" + app.Config.ApiVersion)

	v.GET("/health", app.HealthCheckHandler)

	posts := v.Group("/posts")
	{
		posts.POST("/", app.Store.Posts.CreatePostHandler)

		postsID := posts.Group("/:id")
		{
			postsID.Use(app.Store.Posts.PostContextMiddleware())

			postsID.GET("/", app.Store.Posts.GetPostHandler)
			postsID.PATCH("/", app.Store.Posts.UpdatePostHandler)
			postsID.DELETE("/", app.Store.Posts.DeletePostHandler)
		}
	}

	return r
}
