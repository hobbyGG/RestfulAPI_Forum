package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/contants/code"
	"github.com/hobbyGG/RestfulAPI_forum/contants/contant"
	"github.com/hobbyGG/RestfulAPI_forum/contants/errors"
	"github.com/hobbyGG/RestfulAPI_forum/controllers"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/packages/jwt"
	"go.uber.org/zap"
)

// 写一个jwt认证中间件
// 前端将token按bearer方式存在header传给后端
// 后端解析header得到token
// 解析token后得到用户信息存进ctx中
// 如果没有bearer或token不正确就报错

const (
	authType = "Bearer"
)

func JWTAuth(ctx *gin.Context) {
	// 获取header
	reqHeader := ctx.GetHeader("authorization") //ctx提供的getheader对大小写不敏感，而http包的需要区分
	// getheader的错误反应在string的情况上，空就是错误。这种处理方法简化了代码
	if reqHeader == "" {
		zap.L().Error("ctx.GetHeader error", zap.Error(errors.ErrNeedAuth))
		controllers.ResponseError(ctx, code.NeedAuth)
		ctx.Abort()
		return
	}

	// 判断是否为bearer
	authStr := strings.Split(reqHeader, " ")
	if len(authStr) != 2 {
		zap.L().Error("header Authorization split error", zap.Error(errors.ErrAuthType))
		controllers.ResponseError(ctx, code.InvalidParam)
		ctx.Abort()
		return
	}
	authFmt := authStr[0]
	if authFmt != authType {
		zap.L().Error("authorization type error", zap.Error(errors.ErrAuthType))
		ctx.Abort()
		controllers.ResponseError(ctx, code.AuthType)
		return
	}

	// 解析token
	token := authStr[1]
	cliams, err := jwt.ParseToken(token)
	if err != nil {
		zap.L().Error("jwt.ParseToken error", zap.Error(errors.ErrNeedAuth))
		zap.L().Info("error detail", zap.String("token:", token))
		ctx.Abort()
		controllers.ResponseError(ctx, code.InvalidToken)
		return
	}

	// 检查用户是否在登录状态
	uid := cliams.UID
	uidStr := strconv.Itoa(int(uid))
	redisTokens, err := redis.UserTokenList(uidStr)
	if err != nil {
		zap.L().Error(" redis.UserTokenList error", zap.Error(err))
		ctx.Abort()
		controllers.ResponseError(ctx, code.ServeBusy)
	}
	for _, redisToken := range redisTokens {
		if redisToken == token {
			ctx.Set(contant.StrUID, uid)
			ctx.Set("username", cliams.UserName)
			ctx.Set("token", token)

			// 调用next才会进行下一个handlerfunc，否则将退出
			ctx.Next()
			return
		}
	}
	zap.L().Error("user token error", zap.Error(errors.ErrNotLogin))
	controllers.ResponseError(ctx, code.NotLogin)
	ctx.Abort()
}
