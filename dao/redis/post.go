package redis

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/hobbyGG/RestfulAPI_forum/contants/errors"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"go.uber.org/zap"
)

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

func AddPost(pid, postInfoStr string) error {
	if err := cli.HSet(KeyPostInfoHset, pid, postInfoStr).Err(); err != nil {
		zap.L().Error("cli.HSet error", zap.Error(err))
		return err
	}
	return nil
}

func GetVote(pid, uid int64) (int16, error) {
	key := GetPostVoteKey(pid)
	uidStr := strconv.Itoa(int(uid))
	score, err := cli.ZScore(key, uidStr).Result()
	if err != nil {
		return 0, err
	}
	return int16(score), nil
}

func Vote(pid, uid int64, vote int16) error {
	key := GetPostVoteKey(pid)
	return cli.ZAdd(key, redis.Z{
		Score:  float64(vote),
		Member: uid,
	}).Err()
}

func GetPostVoteKey(pid int64) string {
	pidStr := strconv.Itoa(int(pid))
	return KeyPostVoteZsetPrefix + pidStr + KeyPostVoteZsetSuffix
}
