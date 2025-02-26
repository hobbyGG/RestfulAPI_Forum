package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/contants/code"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/service"
	"go.uber.org/zap"
)

func SignUpHandler(ctx *gin.Context) {
	// 接收用户传来的json数据
	var user models.ParamSignUp
	if err := ctx.ShouldBindJSON(&user); err != nil {
		zap.L().Error("ShouldBindJSON user error", zap.Error(err))
		ResponseError(ctx, code.InvalidParam)
		return
	}

	// 注册服务
	if err := service.SignUp(&user); err != nil {
		zap.L().Error("service.SignUp error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}
	fmt.Printf("%v", user)

	ResponseSuccess(ctx, nil)
}

func LoginHandler(ctx *gin.Context) {
	// 处理登录数据
	var loginParam models.ParamLogin
	if err := ctx.ShouldBindJSON(&loginParam); err != nil {
		zap.L().Error("json bind login error", zap.Error(err))
		ResponseError(ctx, code.InvalidParam)
		return
	}

	// 登录验证信息
	token, err := service.Login(&loginParam)
	if err != nil {
		zap.L().Error("service.Login error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}

	// 多端登录验证
	if err := service.LoginLimit(token); err != nil {
		zap.L().Error("service.LoginLimit error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}

	ResponseSuccess(ctx, token)
}

func LogoutHandler(ctx *gin.Context) {
	// 处理数据
	uid := ctx.GetInt64("uid")
	token := ctx.GetString("token")

	if err := service.Logout(uid, token); err != nil {
		zap.L().Error("service.Logout error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}
	ResponseSuccess(ctx, nil)
}
