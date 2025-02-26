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

	"github.com/hobbyGG/RestfulAPI_forum/dao/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/dao/redis"
	"github.com/hobbyGG/RestfulAPI_forum/log"
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
)

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
	if err := mysql.Init(&settings.Cfg.MysqlCfg); err != nil {
		zap.L().Error("mysql init error", zap.Error(err))
		return
	}

	// 初始化redis
	if err := redis.Init(&settings.Cfg.RedisCfg); err != nil {
		zap.L().Error("redis.Init", zap.Error(err))
		return
	}
	
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

	// 平滑重启
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), shoutDowmTime)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("Shutdown error", zap.Error(err))
	}
	zap.L().Info("shut down succes\nbye~")
}
