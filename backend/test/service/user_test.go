package service_test

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
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/service"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestLoginLimit(t *testing.T) {
	test.Init(t)

	t.Run("login 2 times", func(t *testing.T) {
		loginParam := &models.ParamLogin{
			UserName: "hobbyGG",
			Pwd:      "123",
		}
		token1, err := service.Login(loginParam)
		if err != nil {
			t.Error(err)
			return
		}
		token2, err := service.Login(loginParam)
		if err != nil {
			t.Error(err)
			return
		}

		if err := service.LoginLimit(token1); err != nil {
			t.Error(err)
			return
		}
		if err := service.LoginLimit(token2); err != nil {
			t.Error(err)
			return
		}

		service.Logout(test.TestUID, token1)
		service.Logout(test.TestUID, token2)

	})

	t.Run("log 4 times", func(t *testing.T) {
		tokens := make([]string, 0, 5)
		loginParam := &models.ParamLogin{
			UserName: "hobbyGG",
			Pwd:      "123",
		}

		for i := 0; i < 4; i++ {
			token, err := service.Login(loginParam)
			if err != nil {
				t.Error(err)
				return
			}
			tokens = append(tokens, token)
			if err := service.LoginLimit(token); err != nil {
				t.Error(err)
				return
			}
		}

		defer func(t *testing.T) {
			for _, token := range tokens {
				if err := service.Logout(test.TestUID, token); err != nil {
					t.Error(err)
					return
				}
			}
		}(t)

		// 访问auth路由初始化
		gin.SetMode(gin.DebugMode)
		r := gin.Default()
		w := httptest.NewRecorder()
		body := bytes.NewReader([]byte(""))
		var flag string
		r.POST(test.TestUrl, middleware.JWTAuth, func(ctx *gin.Context) {
			flag = ctx.GetString("token")
		})

		for _, token := range tokens {
			// 如果运行正确，将会有一个token失效访问auth失效
			req, _ := http.NewRequest(http.MethodPost, test.TestUrl, body)
			flag = ""
			v := "Bearer " + token
			req.Header.Set("Authorization", v)

			r.ServeHTTP(w, req)

			if flag != "" {
				// 通过验证的情况
				fmt.Printf("token:%s auth \n", flag)
				continue
			}
			// 没有通过验证
			resData := new(controllers.ResponseParam)
			if err := json.Unmarshal(w.Body.Bytes(), resData); err != nil {
				t.Error(err)
				return
			}

			if resData.Code == code.NotLogin {
				return
			}
			fmt.Println("code:", resData.Code)
		}
		t.Errorf("should have 1 code.NotLogin")
	})

}
