package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	stock "seckill/idl/kitex_gen/stock"
	"seckill/rpc/stock/dao"
	"strconv"
	"time"
)

// StockServiceImpl implements the last service interface defined in the IDL.
type StockServiceImpl struct {
	db      *gorm.DB
	rdb     *redis.Client
	redsync *redsync.Redsync
}

// QueryStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) QueryStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	atoi, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, err
	}
	productStock := dao.RedisSearchStock(s.rdb, atoi)
	if productStock == "" {
		productStock = string(dao.MysqlSearchStock(s.db, atoi).Stock)
	}

	intStock, err := strconv.Atoi(productStock)
	if err != nil {
		return nil, err
	}
	finStock := int32(intStock)
	return &stock.StockResp{Code: 0, Message: "ok", RemainingStock: &finStock}, nil
}

// PreDeductStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) PreDeductStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("stock:%d", productId)

	// 使用 redsync 获取分布式锁，防止多节点同时操作库存
	lockKey := fmt.Sprintf("lock:%s", key)
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second), // 锁 5 秒
		redsync.WithTries(3),              // 最多重试 3 次
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}

	defer func() {
		if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
			// 解锁失败，记录日志
			fmt.Printf("解锁失败: %v\n", err)
		}
	}()

	// Lua 脚本，检查库存是否足够并进行扣减操作
	luaScript := `
        local key = KEYS[1]
        local deduct = tonumber(ARGV[1])
        local stock = tonumber(redis.call("get", key))
        if stock == nil then
            return -1
        end
        if stock >= deduct then
            return redis.call("decrby", key, deduct)
        else
            return -1
        end
    `

	//执行lua脚本
	result, err := s.rdb.Eval(ctx, luaScript, []string{key}, req.Count).Result()
	if err != nil {
		return nil, err
	}
	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "库存不足"}, nil
	}
	remaining := int32(result.(int64))
	return &stock.StockResp{Code: 0, Message: "ok", RemainingStock: &remaining}, nil
}

// RollbackStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) RollbackStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}

// ReserveStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ReserveStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}

// ReleaseStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ReleaseStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}
