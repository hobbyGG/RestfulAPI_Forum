package service

import "errors"

var (
	ErrUserNotExisted = errors.New("用户不存在")
	ErrPwd            = errors.New("用户名或密码错误")
)
