package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt(json web token) 是一种加密传输协议
// 分为头部header，负载payload，签名signature。其中负载中包含声明claims

const (
	expireTime = time.Hour * 7 * 24
	secret     = "7174"
)

var errInvalidToken = errors.New("无效的token")

type MCliams struct {
	UID      int64         `json:"uid"`
	UserName string        `json:"username"`
	ReqTime  time.Duration `json:"reqTime"`
	jwt.StandardClaims
}

func GetToken(uid int64, userName string) (string, error) {
	c := MCliams{
		UID:      uid,
		UserName: userName,
		ReqTime:  time.Duration(time.Now().UnixNano()),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireTime).Unix(), //token过期时间(截止到)
			Issuer:    "rforum",
		},
	}
	// 指定header加密模式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 生成token,传入byte类型，因为底层会对传入值进行改动，使用string会发生类型转换
	return token.SignedString([]byte(secret))
}

// ParseToken 解析token获取用户信息
func ParseToken(tokenStr string) (*MCliams, error) {
	claims := new(MCliams)
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// 这个函数是要找到秘钥
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, errInvalidToken
}
