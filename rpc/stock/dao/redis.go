package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func RedisSearchStock(rdb *redis.Client, id int) string {
	getstring, err := rdb.Get(context.TODO(), strconv.Itoa(id)).Result()
	if err != nil {
		return ""
	}
	return getstring
}
