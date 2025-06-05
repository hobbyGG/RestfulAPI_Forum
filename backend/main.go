package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hobbyGG/RestfulAPI_forum/contants/errors"
	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/log"
	dbsync "github.com/hobbyGG/RestfulAPI_forum/packages/dbSync"
	"github.com/hobbyGG/RestfulAPI_forum/packages/snowflake"
	"github.com/hobbyGG/RestfulAPI_forum/router"
	"github.com/hobbyGG/RestfulAPI_forum/settings"
	"go.uber.org/zap"
)

const (
	setFlag        = "f"
	setUsage       = "项目配置文件路径"
	defaultSetPath = "E:\\Work\\CodeForStudy\\RestfulAPI_Forum\\settings"
	shoutDowmTime  = 5 * time.Second
	retryTimes     = 3
	retryPeriod    = 5 * time.Second
)

// 投票功能重写

func main() {
	// 读取配置文件
	var setPath string
	flag.StringVar(&setPath, setFlag, defaultSetPath, setUsage)
	flag.Parse()
	if err := settings.Init(setPath); err != nil {
		io.WriteString(os.Stdout, err.Error())
		return
	}

	// 初始化zap logger
	if err := log.Init(&settings.Cfg.LogCfg); err != nil {
		io.WriteString(os.Stdout, err.Error())
		return
	}

	// 初始化mysql
	var mysqlInitWrap mysqlInitFunc = mysql.Init
	if err := dbInit(mysqlInitWrap, &settings.Cfg.MysqlCfg); err != nil {
		zap.L().Error("mysql init error", zap.Error(err))
		return
	}
	defer mysql.Close()

	// 初始化redis
	var redisInitWrap redisInitFunc = redis.Init
	if err := dbInit(redisInitWrap, &settings.Cfg.RedisCfg); err != nil {
		zap.L().Error("redis init error", zap.Error(err))
		return
	}
	defer redis.Close()

	// 初始化packages
	if err := snowflake.Init(time.Now().Format("2006-01-02"), 1); err != nil {
		zap.L().Error("snowflake.Init error", zap.Error(err))
		return
	}
	// 初始化路由
	r := router.Init()
	addr := settings.Cfg.Host + ":" + strconv.Itoa(settings.Cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			zap.L().Fatal("ListenAndServe error", zap.Error(err))
		}
	}()

	// 定时刷新mysql分数
	ticker := time.NewTicker(time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := dbsync.Score(); err != nil {
					zap.L().Error("dbsync.Score error", zap.Error(err))
					continue
				}
			case <-done:
				return
			}
		}
	}()

	// 平滑重启
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	done <- true
	ctx, cancel := context.WithTimeout(context.Background(), shoutDowmTime)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("Shutdown error", zap.Error(err))
	}
	zap.L().Info("shut down succes\nbye~")
}

type mysqlInitFunc func(*settings.MysqlCfg) error
type redisInitFunc func(*settings.RedisCfg) error

type InitTypeConstraint interface {
	mysqlInitFunc | redisInitFunc
}

func dbInit[F InitTypeConstraint](f F, c interface{}) error {
	var err error

	switch interface{}(f).(type) {
	case mysqlInitFunc:
		cfg, ok := c.(*settings.MysqlCfg)
		if !ok {
			zap.L().Error("c.(*settings.MysqlCfg) error", zap.Error(errors.ErrCfgType))
			return errors.ErrCfgType
		}

		for i := 0; i < retryTimes; i++ {
			err = nil
			if err = interface{}(f).(mysqlInitFunc)(cfg); err != nil {
				time.Sleep(retryPeriod)
				continue
			}
			break
		}
	case redisInitFunc:
		cfg, ok := c.(*settings.RedisCfg)
		if !ok {
			zap.L().Error("c.(*settings.RedisCfg) error", zap.Error(errors.ErrCfgType))
			return errors.ErrCfgType
		}

		for i := 0; i < retryTimes; i++ {
			err = nil
			if err = interface{}(f).(redisInitFunc)(cfg); err != nil {
				time.Sleep(retryPeriod)
				continue
			}
			break
		}
	default:
		return errors.ErrCfgType
	}
	return err
}
