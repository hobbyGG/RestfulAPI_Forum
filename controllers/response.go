package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/contants/code"
)

type ResponseParam struct {
	Msg  string      `json:"msg"`
	Code code.Code   `json:"code"`
	Data interface{} `json:"data"`
}

func RespSuccess(ctx *gin.Context, data interface{}) {
	resp := &ResponseParam{
		Msg:  code.Success.Msg(),
		Code: code.Success,
		Data: data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func ResponseError(ctx *gin.Context, code code.Code) {
	resp := &ResponseParam{
		Msg:  code.Msg(),
		Code: code,
		Data: nil,
	}
	ctx.JSON(http.StatusOK, resp)
}
