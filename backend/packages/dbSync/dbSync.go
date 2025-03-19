package dbsync

import (
	"strconv"
	"strings"

	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/service"
	"go.uber.org/zap"
)

const KeySocrePattern = "rforum:post:*:vote"

// Score() 将redis中post的分数与mysql进行同步
func Score() error {
	keys, err := redis.GetKeys(KeySocrePattern)
	if err != nil {
		zap.L().Error("redis.GetKeys error", zap.Error(err))
		return err
	}
	for _, key := range keys {
		parts := strings.Split(key, ":")
		pid, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt error", zap.Error(err))
			return err
		}
		if err := service.SyncScoreToMysql(pid); err != nil {
			zap.L().Error("service.SyncScoreToMysql error", zap.Error(err))
			return err
		}
	}
	return nil
}
