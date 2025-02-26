package mysql_test

import (
	"database/sql"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestAddUser(t *testing.T) {
	// 初始化mysql
	test.Init(t)

	// 创建测试用户
	testUser := models.User{
		UID: -1,
	}

	testUser = models.User{
		UID:      -1,
		UserName: "test",
		Pwd:      "test123",
		Email:    "test@test.com",
	}

	// 删除已有的测试用户
	_, err := mysql.DelUserByID(testUser.UID)
	if err != nil {
		if err != sql.ErrNoRows {
			// 如果不是无该用户的错误就是出错了
			t.Errorf(`mysql.DelUserByID error:%v\n`, err)
			return
		}
	}

	if err := mysql.AddUser(&testUser); err != nil {
		t.Error(err)
		return
	}

	// 删除测试用户
	mysql.DelUserByID(-1)
}

func TestGetUserByID(t *testing.T) {
	test.Init(t)

	var id int64 = -1
	user, err := mysql.GetUserByID(id)
	if err != nil {
		t.Errorf("mysql.GetUserByID error:%v\n", err)
		return
	}
	expectedUser := models.User{
		UID:      -1,
		UserName: "test",
		Pwd:      "test123",
		Email:    "test@test.com",
	}
	assert.Equal(t, user, expectedUser)
}

func TestDelUserByID(t *testing.T) {
	test.Init(t)

	var id int64 = -1
	user, err := mysql.DelUserByID(id)
	if err != nil {
		t.Errorf("DelUserByID error:%v\n", err)
		return
	}
	expectedUser := models.User{
		UID:      -1,
		UserName: "test",
		Pwd:      "test123",
		Email:    "test@test.com",
	}
	assert.Equal(t, user, expectedUser)
}

func TestGetUsers(t *testing.T) {
	test.Init(t)

	mysql.DelUserByID(-99)
	mysql.DelUserByID(-100)
	u1, u2 :=
		models.User{
			UID:      -99,
			UserName: "testUser",
			Pwd:      "test123",
			Email:    "test99@test.com",
		},
		models.User{
			UID:      -100,
			UserName: "testUser",
			Pwd:      "test123",
			Email:    "test100@test.com"}
	if err := mysql.AddUser(&u1); err != nil {
		t.Error(err)
		return
	}
	if err := mysql.AddUser(&u2); err != nil {
		t.Error(err)
		return
	}
	defer mysql.DelUserByID(u1.UID)
	defer mysql.DelUserByID(u2.UID)

	testUserName := "testUser"
	users, err := mysql.GetUsers(testUserName)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, users[0], u1)
	assert.Equal(t, users[1], u2)

}
