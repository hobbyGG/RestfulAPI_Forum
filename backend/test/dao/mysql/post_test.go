package mysql_test

import (
	"testing"

	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/test"
)

func TestCreatePost(t *testing.T) {
	// 初始化mysql
	test.Init(t)

	// 声明测试帖子
	post := models.Post{
		PostID:    -1,
		AuthorUID: -1,
		Score:     -1,
		Status:    0,
		ParamCreatePost: models.ParamCreatePost{
			CommID:  1,
			Title:   "test",
			Content: "test",
		},
	}
	// 添加测试帖子
	if err := mysql.CreatePost(&post); err != nil {
		t.Error(err)
		return
	}
	// 删除测试帖子
	mysql.DelPostByID(post.PostID)
}
