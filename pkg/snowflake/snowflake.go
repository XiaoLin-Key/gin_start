package snowflake

import (
	"fmt"
	"gin_start/settings"
	"time"

	sf "github.com/sony/sonyflake"
)

var node *sf.Sonyflake // 注意：sonyflake 的实例类型通常是 *sf.Sonyflake

// Init 初始化雪花算法节点
// startTime: 起始时间，格式如 "2026-02-06"
// machineID: 机器 ID，用于区分分布式环境下的不同节点
func Init(startTime string, machineID uint16) (err error) {
	startTime = settings.Conf.StartTime
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return fmt.Errorf("parse startTime failed: %w", err)
	}

	settings := sf.Settings{
		StartTime: st,
		MachineID: func() (uint16, error) {
			return machineID, nil
		},
	}

	node = sf.NewSonyflake(settings)
	if node == nil {
		return fmt.Errorf("create sonyflake node failed")
	}
	return nil
}

// GenID 生成分布式唯一 ID
func GenID() int64 {
	if node == nil {
		return 0
	}
	id, err := node.NextID()
	if err != nil {
		return 0
	}
	return int64(id)
}
