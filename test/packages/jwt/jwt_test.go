package jwt_test

import (
	"fmt"
	"testing"

	"github.com/hobbyGG/RestfulAPI_forum/packages/jwt"
)

func TestGetToken(t *testing.T) {
	uid := int64(-1)
	UserName := "test"
	token, err := jwt.GetToken(uid, UserName)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjE4MDI2ODcxMTI4NDczNiwidXNlcm5hbWUiOiJob2JieUdHIiwiZXhwIjoxNzQwODkyNzUyLCJpc3MiOiJyZm9ydW0ifQ.aD0oTldMraSnKG4EyNqPX8ZwT3Gd3CLmjHmjw-VDiCY"
	c, err := jwt.ParseToken(token)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(c.UserName)
}
