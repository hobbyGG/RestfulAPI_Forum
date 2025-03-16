package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/hobbyGG/RestfulAPI_forum/settings"
	"go.uber.org/zap"
)

const (
	KeyUserTokenSetPrefix = "rforum:user:"
	KeyUserTokenSetSuffix = ":token"

	KeyPostInfoHset = "rforum:post:info"

	KeyPostVoteZsetPrefix = "rforum:post:"
	KeyPostVoteZsetSuffix = ":vote"
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
