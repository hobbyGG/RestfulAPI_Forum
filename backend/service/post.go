package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/hobbyGG/RestfulAPI_forum/contants/contant"
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

	// 以下应该使用事务操作
	// 将数据存入redis
	if err := redis.UpdatePostRank(post.PostID, post.Score); err != nil {
		zap.L().Error("redis.UpdatePostRank error", zap.Error(err))
		return -2, err
	}

	postInfoStr, err := json.Marshal(post)
	if err != nil {
		zap.L().Error("json.Marshal error", zap.Error(err))
		return -2, err
	}

	pidStr := strconv.Itoa(int(postID))
	if err := redis.AddPostInfo(pidStr, string(postInfoStr)); err != nil {
		zap.L().Error("redis.AddPost error", zap.Error(err))
		return -2, err
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
		if err := redis.AddPostInfo(pidStr, string(postInfoStr)); err != nil {
			zap.L().Error("redis.AddPost error", zap.Error(err))
			return nil, err
		}
	}
	return post, nil
}

// GetPosts 查询某一页指定数量的帖子，并按一定排序返回
func GetPosts(page, size int, sorted string) ([]*models.PostPreview, error) {
	// 判断需要的排行方式
	if sorted == contant.SortedTime {
		return mysql.GetPosts(page, size)
	} else if sorted == contant.SortedScore {
		return redis.GetPosts(page, size)
	}
	return nil, errors.ErrRedisNoKey
}

func PostVote(pid int64, uid int64, vote int16) error {
	// 直接根据投票结果进行修改
	defer SyncScoreToMysql(pid)
	// 先在PostVot修改该用户投票情况
	if vote != 1 && vote != 0 {
		return errors.ErrInvalidParam
	} else if vote == 1 {
		// 直接进行修改
		return redis.Vote(pid, uid, vote)
	} else if vote == 0 {
		return redis.RmVote(pid, uid)
	}

	// 计算新的得分
	score, _ := redis.GetPostScore(pid)

	// 更新postRank表
	redis.UpdatePostRank(pid, score)

	return nil
}

func SyncScoreToMysql(pid int64) error {
	// 根据pid查看对应redis的投票情况
	// 计算投票结果，存入mysql
	votedScore, err := redis.GetPostScore(pid)
	if err != nil {
		zap.L().Error("redis.GetPostScore error", zap.Error(err))
		return err
	}

	mysqlScore, err := mysql.GetPostScore(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostScore error", zap.Error(err))
		return err
	}

	newScore := mysqlScore + votedScore
	if err := mysql.SetPostScore(pid, newScore); err != nil {
		zap.L().Error("mysql.SetPostScore error", zap.Error(err))
		return err
	}
	return nil
}
