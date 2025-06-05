package redis

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/hobbyGG/RestfulAPI_forum/contants/errors"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"go.uber.org/zap"
)

func GetPosts(page, size int) ([]*models.PostPreview, error) {
	key := KeyPostRankZset
	postIDs, err := cli.ZRevRange(key, int64(page), int64(size-1)).Result()
	if err != nil {
		// 表为空
		if err == errors.ErrRedisNil {
			return nil, nil
		}
		// 查询出错
		zap.L().Error("cli.ZRevRange error", zap.Error(err))
		return nil, err
	}
	postsPre := make([]*models.PostPreview, len(postIDs))
	for i, postIDstr := range postIDs {
		postInfo, err := GetPostInfo(postIDstr)
		if err != nil {
			zap.L().Error("GetPostInfo(postIDstr) error", zap.Error(err))
			return nil, err
		}
		json.Unmarshal([]byte(postInfo), &postsPre[i])
	}

	return postsPre, nil
}

func GetPostByID(pid string) (*models.Post, error) {
	postInfoStr, err := cli.HGet(KeyPostInfoHset, pid).Result()
	if err != nil {
		zap.L().Error("cli.HGet error", zap.Error(err))
		return nil, err
	}

	// 如果可以从redis中查到帖子信息
	// 解析json数据
	post := new(models.Post)
	if err := json.Unmarshal([]byte(postInfoStr), post); err != nil {
		zap.L().Error("json.Unmarshal error", zap.Error(err))
		return nil, err
	}

	return post, nil
}

func AddPostInfo(pid, postInfoStr string) error {
	if err := cli.HSet(KeyPostInfoHset, pid, postInfoStr).Err(); err != nil {
		zap.L().Error("cli.HSet error", zap.Error(err))
		return err
	}
	return nil
}

func GetPostInfo(pid string) (string, error) {
	return cli.HGet(KeyPostInfoHset, pid).Result()
}

func Vote(pid, uid int64, vote int16) error {
	key := GetPostVoteKey(pid)
	return cli.SAdd(key, uid).Err()
}

func RmVote(pid, uid int64) error {
	key := GetPostVoteKey(pid)
	return cli.SRem(key, uid).Err()
}

func GetPostVoteKey(pid int64) string {
	pidStr := strconv.Itoa(int(pid))
	return KeyPostVoteSetPrefix + pidStr + KeyPostVoteSetSuffix
}

// GetPostScore 根据帖子id返回该帖子的分数
func GetPostScore(pid int64) (int64, error) {
	key := KeyPostRankZset
	pidStr := strconv.Itoa(int(pid))
	score, err := cli.ZScore(key, pidStr).Result()
	if err != nil {
		zap.L().Error("cli.ZScore(key, pidStr) error", zap.Error(err))
		return -1, err
	}

	return int64(score), err
}

func GetPostVoteLen(pid int64) int64 {
	key := GetPostVoteKey(pid)
	l, _ := cli.SCard(key).Result()
	return l
}

func UpdatePostRank(pid, score int64) error {
	key := KeyPostRankZset

	return cli.ZAdd(key, redis.Z{
		Member: pid,
		Score:  float64(score),
	}).Err()
}
