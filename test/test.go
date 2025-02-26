package test

import (
	"testing"
	"time"

	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/log"
	"github.com/hobbyGG/RestfulAPI_forum/packages/snowflake"
	"github.com/hobbyGG/RestfulAPI_forum/settings"
)

const (
	defaultSetPath = "E:\\Work\\CodeForStudy\\RestfulAPI_Forum\\settings"
	TestUrl = "/iusadhfglianfdlgahcv/test"
	TestUID        = 180268711284736
)

func Init(t *testing.T) {
	t.Helper()
	// 初始化mysql
	if err := settings.Init(defaultSetPath); err != nil {
		t.Error(err)
		return
	}
	if err := mysql.Init(&settings.Cfg.MysqlCfg); err != nil {
		t.Error(err)
		return
	}
	if err := redis.Init(&settings.Cfg.RedisCfg); err != nil {
		t.Error(err)
		return
	}
	if err := log.Init(&settings.Cfg.LogCfg); err != nil {
		t.Error(err)
		return
	}
	if err := snowflake.Init(time.Now().Format("2006-01-02"), 1); err != nil {
		t.Error(err)
		return
	}
}
