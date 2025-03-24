package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/hobbyGG/RestfulAPI_forum/settings"
	"go.uber.org/zap"
)

const (
	// 以登录的用户列表
	KeyUserTokenSetPrefix = "rforum:user:"
	KeyUserTokenSetSuffix = ":token"

	// 帖子信息
	KeyPostInfoHset = "rforum:post:info"

	// 每个帖子的投票情况
	KeyPostVoteZsetPrefix = "rforum:post:"
	KeyPostVoteZsetSuffix = ":vote"

	// 帖子排行
	KeyPostRankZset = "rforum:post:rank"
)

var cli *redis.Client

type InitFuncType func(*settings.RedisCfg) error

func Init(cfg *settings.RedisCfg) error {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	cli = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Pwd,
		DB:       cfg.DBnum,
	})

	_, err := cli.Ping().Result()
	if err != nil {
		zap.L().Error("cli.Ping() error", zap.Error(err))
		return err
	}
	return nil
}

func Close() {
	if err := cli.Save().Err(); err != nil {
		zap.L().Error("cli.Save() error", zap.Error(err))
		// 持久化失败，应该启用其他保存方案，比如存入mysql
		return
	}
	if err := cli.Close(); err != nil {
		zap.L().Error("cli.Close error", zap.Error(err))
	}
}

func GetKeys(pattern string) ([]string, error) {
	var keys []string
	var cursor uint64
	for {
		ks, c, err := cli.Scan(cursor, pattern, 100).Result()
		if err != nil {
			zap.L().Error("cli.Scan error", zap.Error(err))
			return nil, err
		}
		keys = append(keys, ks...)
		cursor = c
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}
