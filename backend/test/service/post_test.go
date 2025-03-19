package service_test

import (
	"database/sql"
	"testing"

	"github.com/hobbyGG/RestfulAPI_forum/service"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestGetPost(t *testing.T) {
	test.Init(t)

	t.Run("pid is wrong", func(t *testing.T) {
		pidStr := "-1"
		_, err := service.GetPost(pidStr)
		if err == sql.ErrNoRows {
			return
		}
		t.Error("post should not existed")
	})

	t.Run("pid is correct", func(t *testing.T) {
		pidStr := "213826083491840"
		post, err := service.GetPost(pidStr)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(post)
	})
}
