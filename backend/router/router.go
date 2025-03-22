package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/contants/contant"
	"github.com/hobbyGG/RestfulAPI_forum/controllers"
	"github.com/hobbyGG/RestfulAPI_forum/middleware"
	"github.com/hobbyGG/RestfulAPI_forum/service"
)

func Init() *gin.Engine {
	r := gin.Default()

	// 配置 CORS 中间件
	config := cors.DefaultConfig()
	// 允许所有源访问，实际生产环境建议指定具体的源
	config.AllowAllOrigins = true
	// 允许的请求方法
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	// 允许的请求头
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	// 允许携带凭证（如 Cookie）
	config.AllowCredentials = true
	// 预检请求的缓存时间
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	r.OPTIONS("/*path", func(ctx *gin.Context) {
		// 设置允许的源，生产环境建议指定具体源，而非使用 *
		ctx.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		// 设置允许的请求头
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")
		// 允许携带凭证
		ctx.Header("Access-Control-Allow-Credentials", "true")
		// 设置预检请求的缓存时间
		ctx.Header("Access-Control-Max-Age", "43200") // 缓存 12 小时
		// 返回 204 状态码
		ctx.Status(http.StatusOK)

	})

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
		api.POST("/loginAuth", loginAuthHandler)
	}

	return r
}

func pingHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func loginAuthHandler(ctx *gin.Context) {
	uid := ctx.GetInt64(contant.StrUID)
	username, _ := service.GetUsername(uid)
	controllers.ResponseSuccess(ctx, username)
}
