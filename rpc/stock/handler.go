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

	// Redis 中预留库存的 key
	reservedKey := fmt.Sprintf("reserved:%d", productId)
	// Redlock 锁 key
	lockKey := fmt.Sprintf("lock:%s", reservedKey)

	// 创建 Redlock 互斥锁
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second), // 锁 5 秒
		redsync.WithTries(3),              // 最多重试 3 次
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer func() {
		if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
			fmt.Printf("解锁失败: %v\n", err)
		}
	}()

	// Lua 脚本：检查并扣减 `reserved:product_id` 的数量
	luaScript := `
        local reservedKey = KEYS[1]
        local deduct = tonumber(ARGV[1])
        local reservedStock = tonumber(redis.call("get", reservedKey))
        if reservedStock == nil or reservedStock < deduct then
            return -1
        end
        return redis.call("decrby", reservedKey, deduct)
    `

	// 执行 Lua 脚本，尝试扣减预订的库存
	result, err := s.rdb.Eval(ctx, luaScript, []string{reservedKey}, req.Count).Result()
	if err != nil {
		return nil, err
	}
	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "预订库存不足"}, nil
	}
	remaining := int32(result.(int64))

	// 更新 MySQL 中的库存
	tx := s.db.Begin()
	if err := tx.Exec("UPDATE product_stocks SET stock = stock - ? WHERE product_id = ?", req.Count, productId).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新 MySQL 库存失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	return &stock.StockResp{Code: 0, Message: "确认扣减成功", RemainingStock: &remaining}, nil
}

// RollbackStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) RollbackStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("stock:%d", productId)
	lockKey := fmt.Sprintf("lock:%s", key)
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second), // 锁 5 秒
		redsync.WithTries(3),              // 最多重试 3 次
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer mutex.UnlockContext(ctx)
	// Lua 脚本操作：回滚 Redis 中的库存
	luaScript := `
        local key = KEYS[1]
        local rollback = tonumber(ARGV[1])
        local stock = tonumber(redis.call("get", key))
        if stock == nil then
            return -1
        end
        return redis.call("incrby", key, rollback)
    `

	result, err := s.rdb.Eval(ctx, luaScript, []string{key}, req.Count).Result()
	if err != nil {
		return nil, fmt.Errorf("Redis 执行失败: %v", err)
	}

	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "库存回滚失败"}, nil
	}

	// 更新 MySQL 库存
	tx := s.db.Begin()

	// 使用 GORM 更新 MySQL 库存
	if err := tx.Exec("UPDATE product_stocks SET stock = stock + ? WHERE product_id = ?", req.Count, productId).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新 MySQL 库存失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	remaining := int32(result.(int64))

	return &stock.StockResp{Code: 0, Message: "库存回滚成功", RemainingStock: &remaining}, nil
}

// ReserveStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ReserveStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, err
	}

	// 使用 Redlock 获取分布式锁
	key := fmt.Sprintf("stock:%d", productId)
	lockKey := fmt.Sprintf("lock:%s", key)
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second),
		redsync.WithTries(3),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer mutex.UnlockContext(ctx)

	// 在 Redis 中预占库存
	redisKey := fmt.Sprintf("reserved:%d", productId)
	stockCount, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取 Redis 库存失败: %v", err)
	}

	remainingStock, _ := strconv.Atoi(stockCount)
	if remainingStock <= 0 {
		return &stock.StockResp{Code: 1, Message: "库存不足"}, nil
	}

	// 使用 Lua 脚本确保库存的预占是原子操作
	luaScript := `
        local key = KEYS[1]
        local stock = tonumber(redis.call("get", key))
        if stock == nil or stock <= 0 then
            return -1
        end
        redis.call("decrby", key, 1)  -- 扣减库存
        return stock - 1
    `
	result, err := s.rdb.Eval(ctx, luaScript, []string{redisKey}).Result()
	if err != nil {
		return nil, fmt.Errorf("执行 Redis Lua 脚本失败: %v", err)
	}

	if result == int64(-1) {
		return &stock.StockResp{Code: 1, Message: "库存不足，预占失败"}, nil
	}

	remaining := int32(result.(int64))
	return &stock.StockResp{Code: 0, Message: "库存预占成功", RemainingStock: &remaining}, nil
}

// ReleaseStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ReleaseStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, err
	}

	// 使用 Redlock 获取分布式锁
	key := fmt.Sprintf("stock:%d", productId)
	lockKey := fmt.Sprintf("lock:%s", key)
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second),
		redsync.WithTries(3),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer mutex.UnlockContext(ctx)

	// Redis Lua 脚本：检查预占库存并释放库存
	luaScript := `
        local stockKey = KEYS[1]
        local reservedKey = KEYS[2]
        local productId = ARGV[1]
        
        local reservedStock = tonumber(redis.call("GET", reservedKey))
        if reservedStock == nil or reservedStock <= 0 then
            return {err="没有预占库存，无法释放"}
        end
        
        -- 执行释放库存操作
        redis.call("DECR", reservedKey)
        redis.call("INCR", stockKey)
        
        return {ok="库存释放成功"}
    `

	reservedKey := fmt.Sprintf("reserved:%d", productId)
	stockKey := fmt.Sprintf("stock:%d", productId)

	// 执行 Lua 脚本
	result, err := s.rdb.Eval(ctx, luaScript, []string{stockKey, reservedKey}, productId).Result()
	if err != nil {
		return nil, fmt.Errorf("执行 Redis Lua 脚本失败: %v", err)
	}

	// 处理返回结果
	resultMap, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("解析 Redis Lua 脚本结果失败")
	}

	// 解析返回的结果
	if errMsg, ok := resultMap[0].(string); ok && errMsg == "没有预占库存，无法释放" {
		return &stock.StockResp{Code: 1, Message: "没有预占库存，无法释放"}, nil
	}

	// 如果脚本执行成功，返回成功信息
	return &stock.StockResp{Code: 0, Message: "库存释放成功"}, nil
}
