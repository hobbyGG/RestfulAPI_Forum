package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/contants/code"
	"github.com/hobbyGG/RestfulAPI_forum/controllers"
	"github.com/hobbyGG/RestfulAPI_forum/middleware"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestCreatePost(t *testing.T) {
	// 初始化mysql
	test.Init(t)

	// 初始化路由
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/post"
	r.POST(url, middleware.JWTAuth, controllers.CreatePostHandler)

	bodyJSON := []byte(`
	{
    "title": "支持中文帖子",
    "content": "可以支持提交创建中文的帖子，还可以选择分区，默认为分区1"
	}`)
	body := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjE4MDI2ODcxMTI4NDczNiwidXNlcm5hbWUiOiJob2JieUdHIiwiZXhwIjoxNzQwOTgzNTA0LCJpc3MiOiJyZm9ydW0ifQ.MYOAiWCqbGOvquLmClT4RL0PgzW3kmijwYdDl8fy6Oc")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	var resData controllers.ResponseParam
	if err := json.Unmarshal(w.Body.Bytes(), &resData); err != nil {
		t.Error(err)
		return
	}

	if resData.Code != code.Success {
		t.Errorf("want 0, but get %v", resData.Code)
		return
	}
}
