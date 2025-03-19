package redis_test

import (
	"testing"

	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestUserToken(t *testing.T) {
	// 初始化redis
	test.Init(t)
	uid := "-1"
	token := "test"
	if err := redis.AddUserToken(uid, token); err != nil {
		t.Error(err)
		return
	}
	if err := redis.SubUserToken(uid); err != nil {
		t.Error(err)
		return
	}
}

func TestUserLoginNum(t *testing.T) {
	test.Init(t)

	uid := "-1"
	token1 := "testToken1"
	token2 := "testToken2"

	if err := redis.AddUserToken(uid, token1); err != nil {
		t.Error(err)
		return
	}
	if err := redis.AddUserToken(uid, token2); err != nil {
		t.Error(err)
		return
	}

	n, err := redis.UserTokenNum(uid)
	if err != nil {
		t.Error(err)
		return
	}

	if n != 2 {
		t.Errorf("want 2, but get %d", n)
		return
	}

	for i := 0; i < 2; i++ {
		redis.SubUserToken(uid)
	}
}
