package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hobbyGG/RestfulAPI_forum/settings"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

var db *sqlx.DB

type InitFuncType func(*settings.MysqlCfg) error

func Init(cfg *settings.MysqlCfg) error {
	// 使用sqlx的connect连接数据库，connect是open和ping的集成
	userName := cfg.UserName
	pwd := cfg.Pwd
	host := cfg.Host
	port := cfg.Port
	dbName := cfg.DBname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", userName, pwd, host, port, dbName)
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("sqlx.Connect error", zap.Error(err))
		return err
	}
	return nil
}

func Close() {
	db.Close()
}
