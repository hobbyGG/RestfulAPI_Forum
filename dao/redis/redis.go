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

// uid token
