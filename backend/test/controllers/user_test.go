package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/hobbyGG/RestfulAPI_forum/controllers"
	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/packages/jwt"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

const defaultSetPath = "E:\\Work\\CodeForStudy\\RestfulAPI_Forum\\settings"

func TestSignUpHandler(t *testing.T) {

	// 初始化路由
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/signup"
	r.POST(url, controllers.SignUpHandler)

	// 新建请求
	newUser := models.ParamSignUp{
		UserName: "test",
		Pwd:      "123123",
		RePwd:    "123123",
		Email:    "test@test.com",
	}
	bodyStr, _ := json.Marshal(&newUser)
	buf := bytes.NewReader(bodyStr)
	req, _ := http.NewRequest(http.MethodPost, url, buf)

	// 新建响应
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	// 解析的到的数据
	res := new(controllers.ResponseParam)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Errorf("unmarshal error:%v, w.bodyStr:%s", err, w.Body.Bytes())
		return
	}
	if res.Code != 2 {
		t.Error("want 2, got ", res.Code)
		return
	}
}

func TestLoinHandelr(t *testing.T) {
	// 初始化路由
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/login"
	r.POST(url, controllers.LoginHandler)
	// 初始化数据库
	test.Init(t)

	// 注册一个用户
	testUser := models.User{
		UID:      -1,
		UserName: "test",
		Pwd:      "test123",
		Email:    "test@test.com",
	}
	mysql.AddUser(&testUser)
	defer mysql.DelUserByID(testUser.UID)

	// 发送登录请求
	loginParam := models.ParamLogin{
		UserName: "test",
		Pwd:      "test123",
	}
	body, err := json.Marshal(&loginParam)
	if err != nil {
		t.Error(err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Error(err)
		return
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 解析body
	res := controllers.ResponseParam{}
	json.Unmarshal(w.Body.Bytes(), &res)
	// 验证响应token
	c, err := jwt.ParseToken(res.Data.(string))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, c.UID, testUser.UID)
}
