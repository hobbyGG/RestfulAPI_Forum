package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/controllers"
	"github.com/hobbyGG/RestfulAPI_forum/middleware"
)

func Init() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api", middleware.RateLimit(time.Second*2, 5))
	{
		api.POST("/signup", controllers.SignUpHandler)
		api.POST("/login", controllers.LoginHandler)

		api.Use(middleware.JWTAuth)

		api.POST("/logout", controllers.LogoutHandler)
		api.POST("/post", controllers.CreatePostHandler)
		api.POST("/vote", controllers.PostVoteHandler)

		api.GET("/post/:id", controllers.GetPostHander)
		api.GET("/posts", controllers.GetPostsHander)

		api.GET("/ping", pingHandler)
	}

	return r
}

func pingHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
