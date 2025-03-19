package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/middleware"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestJWTAuth(t *testing.T) {
	test.Init(t)
	// 初始化路由
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/test/auth"
	var uid int64
	r.POST(url, middleware.JWTAuth, func(ctx *gin.Context) {
		uid = ctx.GetInt64("uid")
	})

	// 建立redis登录名单
	redis.AddUserToken("180268711284736", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjE4MDI2ODcxMTI4NDczNiwidXNlcm5hbWUiOiJob2JieUdHIiwiZXhwIjoxNzQwOTgzNTA0LCJpc3MiOiJyZm9ydW0ifQ.MYOAiWCqbGOvquLmClT4RL0PgzW3kmijwYdDl8fy6Oc")
	defer redis.SubUserToken("180268711284736")
	// 新建请求
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjE4MDI2ODcxMTI4NDczNiwidXNlcm5hbWUiOiJob2JieUdHIiwiZXhwIjoxNzQwOTgzNTA0LCJpc3MiOiJyZm9ydW0ifQ.MYOAiWCqbGOvquLmClT4RL0PgzW3kmijwYdDl8fy6Oc")

	// 新建响应
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if uid == 0 {
		t.Errorf("auth fail, uid is null")
		return
	}
	fmt.Printf("uid:%d", uid)
}
