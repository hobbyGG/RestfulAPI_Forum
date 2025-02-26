package service

import (
	"fmt"
	"strconv"

	"github.com/hobbyGG/RestfulAPI_forum/contants/contant"
	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/packages/jwt"
	"github.com/hobbyGG/RestfulAPI_forum/packages/snowflake"
	"go.uber.org/zap"
)

func SignUp(user *models.ParamSignUp) error {
	// 新建用户
	uid, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID error", zap.Error(err))
		return err
	}
	newUser := &models.User{
		UID:      uid,
		UserName: user.UserName,
		Pwd:      user.Pwd,
		Email:    user.Email,
	}
	return mysql.AddUser(newUser)
}

// Login 正确运行的情况下回返回token的值与一个错误信息
func Login(loginParam *models.ParamLogin) (string, error) {
	// 检查用户是否存在(用户名可以重复)
	users, err := mysql.GetUsers(loginParam.UserName)
	if err != nil {
		zap.L().Error("mysql.GetUsers", zap.Error(err))
		return "", err
	}
	if len(users) == 0 {
		return "", ErrUserNotExisted
	}

	// 检查用户名与密码是否匹配
	pwd := loginParam.Pwd
	for _, user := range users {
		if pwd == user.Pwd {
			token, err := jwt.GetToken(user.UID, user.UserName)
			if err != nil {
				return "", err
			}
			return token, nil
		}
	}
	return "", ErrPwd
}

func LoginLimit(token string) error {
	// 根据token判断用户登录了几个端
	// 获取用户的uid
	cliams, err := jwt.ParseToken(token)
	if err != nil {
		zap.L().Error("jwt.ParseToken error", zap.Error(err))
		return err
	}
	uid := cliams.UID
	uidStr := strconv.Itoa(int(uid))

	// 通过用户uid查看用户已经登录了几次
	n, err := redis.UserTokenNum(uidStr)
	if err != nil {
		zap.L().Error("redis.UserLoginNum error", zap.Error(err))
		return err
	}

	if n >= contant.MaxUserLogin {
		// 超过最大用户登录次数，随机删去一个用户的token，改为该用户的token
		if err := redis.SubUserToken(uidStr); err != nil {
			zap.L().Error("redis.SubUserToken error", zap.Error(err))
			return err
		}
		if err := redis.AddUserToken(uidStr, token); err != nil {
			zap.L().Error("redis.AddUserToken error", zap.Error(err))
			return err
		}
		zap.L().Info(fmt.Sprintf("user %s login num over %d", uidStr, contant.MaxUserLogin))
		return nil
	}

	// 记录登录的用户
	if err := redis.AddUserToken(uidStr, token); err != nil {
		zap.L().Error("n < max, but redis.AddUserToken error", zap.Error(err))
		return err
	}

	return nil
}

func Logout(uid int64, token string) error {
	uidStr := strconv.Itoa(int(uid))
	return redis.Logout(uidStr, token)
}
