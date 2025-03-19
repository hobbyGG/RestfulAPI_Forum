package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimit(interval time.Duration, cap int64) func(ctx *gin.Context) {
	bucket := ratelimit.NewBucket(interval, cap)
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) == 0 {
			// takeAvailable就是取几个令牌，返回0则是需要等待
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "点击的用户太多了，稍后重试",
			})
			ctx.Abort()
			return
		}
		// 成功取到令牌
		ctx.Next()
	}
}
