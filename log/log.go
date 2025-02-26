package log

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hobbyGG/RestfulAPI_forum/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logDirPath = "E:\\Work\\CodeForStudy\\RestfulAPI_Forum\\log\\"
	debugMode  = "debug"
)

func Init(cfg *settings.LogCfg) error {
	// 初始化写同步器
	var buf bytes.Buffer
	buf.WriteString(logDirPath)
	buf.WriteString(cfg.FileName)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   buf.String(),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackup,
		MaxAge:     cfg.MaxAge,
	}
	lumberWriteSyncer := zapcore.AddSync(lumberjackLogger)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	// 初始化编码器 NewProduction()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // 修改时间表示
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	coreFile := zapcore.NewCore(encoder, lumberWriteSyncer, zapcore.DebugLevel)
	coreConsole := zapcore.NewCore(encoder, consoleWriteSyncer, zapcore.DebugLevel)

	lg := new(zap.Logger)
	if settings.Cfg.Mode == debugMode {
		coreTee := zapcore.NewTee(coreFile, coreConsole)
		lg = zap.New(coreTee)
		fmt.Println("project running in debug mode")
	} else {
		lg = zap.New(coreFile)
		fmt.Println("project running in dev mode")
	}

	// 注册全局可调用的日志
	zap.ReplaceGlobals(lg)
	return nil
}
