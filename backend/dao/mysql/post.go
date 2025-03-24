package mysql

import (
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"go.uber.org/zap"
)

func CreatePost(post *models.Post) error {
	sqlStr := `
	insert into post(postID, authorUID, score, status, commID, title, content)
	values(?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, post.PostID, post.AuthorUID, post.Score, post.Status, post.CommID, post.Title, post.Content)
	if err != nil {
		zap.L().Error("db.Exec insert post error", zap.Error(err))
		return err
	}
	return nil
}

func DelPostByID(postID int64) error {
	sqlStr := `
	delete from post where postID=?`
	if _, err := db.Exec(sqlStr, postID); err != nil {
		zap.L().Error("delete post error", zap.Error(err))
		return err
	}
	return nil
}

func GetPostByID(pid int64) (*models.Post, error) {
	sqlStr := `
	select postID, authorUID, score, status, commID, title, content, create_time from post
	where postID = ?`

	post := new(models.Post)
	if err := db.Get(post, sqlStr, pid); err != nil {
		zap.L().Error("db.Get error", zap.Error(err))
		return nil, err
	}

	return post, nil
}

func GetPosts(page, size int) ([]*models.PostPreview, error) {
	sqlStr := `
	select postID, title, content, score, commID, status, create_time
	from post
	order by create_time desc
	limit ? offset ?`
	offset := page * size

	posts := make([]*models.PostPreview, 0, 10)
	if err := db.Select(&posts, sqlStr, size, offset); err != nil {
		zap.L().Error("db.Select error", zap.Error(err))
		return nil, err
	}

	for _, post := range posts {
		p, _ := GetPostByID(post.PostID)
		uid := p.AuthorUID
		user, err := GetUserByID(uid)
		if err != nil {
			zap.L().Error("GetUserByID error", zap.Error(err))
			return nil, err
		}
		post.AuthorName = user.UserName
	}
	return posts, nil
}

func GetPostScore(pid int64) (int64, error) {
	sqlStr := `select score from post where postID = ?`
	score := new(int64)
	if err := db.Get(score, sqlStr, pid); err != nil {
		zap.L().Error("db.Get error", zap.Error(err))
		return 0, err
	}

	return *score, nil
}

func SetPostScore(pid, score int64) error {
	sqlStr := `update post set score = ? where postID = ?`
	if _, err := db.Exec(sqlStr, score, pid); err != nil {
		zap.L().Error("db.Exec error", zap.Error(err))
		return err
	}
	return nil
}
