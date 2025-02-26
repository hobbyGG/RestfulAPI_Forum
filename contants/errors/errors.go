package errors

import (
	"errors"

	"github.com/go-redis/redis"
)

var (
	ErrUserNotExisted = errors.New("用户不存在")
	ErrPwd            = errors.New("用户名或密码错误")
	ErrNotLogin       = errors.New("用户未登录")
	ErrNeedAuth       = errors.New("需要token")
	ErrInvalidHeader  = errors.New("header格式错误")
	ErrAuthType       = errors.New("错误的auth方式")
	ErrTokenExpired   = errors.New("用户token已过期")
	ErrNeedPID        = errors.New("需要帖子的id")

	ErrRedisNil = redis.Nil
)
