package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/hobbyGG/RestfulAPI_forum/contants/errors"
	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/packages/snowflake"
	"go.uber.org/zap"
)

const defauleStatus = 0

func CreatePost(postParam *models.ParamCreatePost, uid int64) (int64, error) {
	// 补全post信息
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID error", zap.Error(err))
		return -2, err
	}
	score := time.Now().Unix()
	post := &models.Post{
		PostID:          postID,
		AuthorUID:       uid,
		Score:           score,
		Status:          defauleStatus,
		ParamCreatePost: *postParam,
	}

	// 将数据存入mysql
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost", zap.Error(err))
		return -2, err
	}
	return postID, nil
}

func GetPost(pidStr string) (*models.Post, error) {

	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt error", zap.Error(err))
		return nil, err
	}

	post, err := redis.GetPostByID(pidStr)
	if err != nil {
		if err != errors.ErrRedisNil {
			zap.L().Error("redis.GetPostByID error", zap.Error(err))
			return nil, err
		}

		// 处理redis中没有存储该帖子的问题
		// 先从mysql中查找
		post, err = mysql.GetPostByID(pid)
		if err != nil {
			// mysql中也没有则错误
			zap.L().Error("mysql.GetPostByID error", zap.Error(err))
			return nil, err
		}

		// 在mysql中找到后记录redis表
		postInfoStr, err := json.Marshal(post)
		if err != nil {
			zap.L().Error("json.Marshal error", zap.Error(err))
			return nil, err
		}
		if err := redis.AddPost(pidStr, string(postInfoStr)); err != nil {
			zap.L().Error("redis.AddPost error", zap.Error(err))
			return nil, err
		}
	}
	return post, nil
}

// GetPosts 查询某一页指定数量的帖子，并按一定排序返回
func GetPosts(page, size int, sorted string) ([]*models.PostPreview, error) {
	return mysql.GetPosts(page, size, sorted)
}

func PostVote(pid int64, uid int64, vote int16) error {
	// 处理用户投票的情况有
	// 0 	->  1/-1
	// 1/-1 -> -1/1
	voted, err := redis.GetVote(pid, uid)
	if err != nil {
		if err == errors.ErrRedisNil {
			// 用户还没投过票或者帖子还没投过票
			err = redis.Vote(pid, uid, vote)
			return err
		}
		zap.L().Error("redis.GetVote error", zap.Error(err))
		return err
	}

	newVote := voted - vote
	if err := redis.Vote(pid, uid, newVote); err != nil {
		zap.L().Error("redis.Vote error", zap.Error(err))
		return err
	}

	return nil
}

// vote情况在redis的zset中维护
// zset表名为post的id，成员为投票的用户
// 如果vote的帖子不在redis中再从mysql读入
