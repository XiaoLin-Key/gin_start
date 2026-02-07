package redis

import (
	"context"
	"fmt"
	"gin_start/settings"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background() // 定义全局上下文，方便内部调用
)

func InitRedis() error {
	conf := settings.Conf.Redis

	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
	})
	// v9 必须传 Ctx
	return Rdb.Ping(Ctx).Err()
}
