package settings

import (
	"io"
	"os"

	"github.com/spf13/viper"
)

type AppCfg struct {
	// 由于要用到映射，所以开头必须大写，否则外部无法读取
	PrjName  string `mapstructure:"prj_name"`
	Version  string `mapstructure:"version"`
	Mode     string `mapstructure:"mode"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"` //项目启动的端口
	LogCfg   `mapstructure:"log"`
	MysqlCfg `mapstructure:"mysql"`
	RedisCfg `mapstructure:"redis"`
}

type LogCfg struct {
	FileName  string `mapstructure:"file_name"`
	Level     string `mapstructure:"level"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
	MaxAge    int    `mapstructure:"max_age"`
}

type MysqlCfg struct {
	UserName string `mapstructure:"user_name"`
	Pwd      string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DBname   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
}

type RedisCfg struct {
	Host  string `mapstructure:"host"`
	Port  string `mapstructure:"port"`
	Pwd   string `mapstructure:"pwd"`
	DBnum int    `mapstructure:"dbnum"`
}

var Cfg AppCfg

func Init(filePath string) error {
	// 添加文件位置
	viper.AddConfigPath(filePath)
	// 读取文件
	if err := viper.ReadInConfig(); err != nil {
		io.WriteString(os.Stdout, err.Error())
		return err
	}

	// 使用unmarshal将参数绑定到结构体中
	if err := viper.Unmarshal(&Cfg); err != nil {
		io.WriteString(os.Stdout, err.Error())
		return err
	}
	return nil
}
