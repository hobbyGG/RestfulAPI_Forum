package snowflake

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *snowflake.Node

var errNilPointer = errors.New("所用的指针为空")

// 配置开始时间与机器id
func Init(startTime string, machineID int64) error {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		zap.L().Error("time.Parse error", zap.Error(err))
		return err
	}
	snowflake.Epoch = st.UnixNano() / 1000000 //纳秒转为毫秒
	node, err = snowflake.NewNode(machineID)
	return err
}

func GetID() (int64, error) {
	if node == nil {
		return 0, errNilPointer
	}
	return node.Generate().Int64(), nil
}
