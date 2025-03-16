package initialize

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"seckill/rpc/stock/global"
)

func InitRedisDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Config.RedisConfig.Addr,
		Password: global.Config.RedisConfig.Password,
		DB:       global.Config.RedisConfig.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}

// 初始化 Redis 同步锁
func InitRedisSync() *redsync.Redsync {
	// 创建 Redis 客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"localhost:6379"}, // 你的 Redis 地址
		Password: "",                         // Redis 密码（如果有）
		DB:       0,                          // 选择的数据库
	})

	// 使用 redigo 适配 go-redis
	pool := goredis.NewPool(rdb)

	// 创建 Redsync 实例
	rs := redsync.New(pool)

	return rs
}
