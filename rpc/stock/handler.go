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
	// 将产品ID转换为整数
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("产品ID格式错误: %v", err)
	}

	// 统一使用 key "stock:productId" 存储商品库存
	// 从 Redis 查询库存（dao 层内部应使用与这里相同的 key 规则）
	productStock := dao.RedisSearchStock(s.rdb, productId)
	// 如果 Redis 中未命中，则回退到 MySQL 查询
	if productStock == "" {
		// 应使用 strconv.Itoa
		productStock = strconv.Itoa(int(dao.MysqlSearchStock(s.db, productId).Stock))
	}

	// 将库存字符串转换为整数
	intStock, err := strconv.Atoi(productStock)
	if err != nil {
		return nil, fmt.Errorf("库存转换错误: %v", err)
	}
	finStock := int32(intStock)
	return &stock.StockResp{Code: 0, Message: "ok", RemainingStock: &finStock}, nil
}

// PreDeductStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) PreDeductStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("产品ID格式错误: %v", err)
	}

	reservedKey := fmt.Sprintf("reserved:%d", productId)
	lockKey := fmt.Sprintf("lock:%s", reservedKey)

	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second),
		redsync.WithTries(3),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer func() {
		if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
			fmt.Printf("解锁失败: %v\n", err)
		}
	}()

	// 先更新 MySQL
	tx := s.db.Begin()
	if err := tx.Exec("UPDATE product_stocks SET stock = stock - ? WHERE product_id = ?", req.Count, productId).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新 MySQL 库存失败: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	// 然后扣减 Redis reserved
	luaScript := `
        local key = KEYS[1]
        local deduct = tonumber(ARGV[1])
        local current = tonumber(redis.call("GET", key))
        if current == nil or current < deduct then
            return -1
        end
        local newStock = current - deduct
        if newStock < 0 then
            return -1
        end
        redis.call("SET", key, newStock)
        return newStock
    `

	result, err := s.rdb.Eval(ctx, luaScript, []string{reservedKey}, req.Count).Result()
	if err != nil {
		return nil, fmt.Errorf("执行 Redis Lua 脚本失败: %v", err)
	}
	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "预订库存不足"}, nil
	}
	remaining := int32(result.(int64))

	return &stock.StockResp{Code: 0, Message: "预扣库存成功", RemainingStock: &remaining}, nil
}

// RollbackStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) RollbackStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// 将产品ID转换为整数
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("产品ID格式错误: %v", err)
	}

	// 统一使用 key "stock:%d" 存储商品库存
	redisKey := fmt.Sprintf("stock:%d", productId)
	// 构造分布式锁的 key
	lockKey := fmt.Sprintf("lock:%s", redisKey)
	// 获取锁，确保并发环境下库存回滚操作的原子性
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second),
		redsync.WithTries(3),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer mutex.UnlockContext(ctx)

	// Lua 脚本：原子地将库存加回 Redis 中
	luaScript := `
        local key = KEYS[1]
        local rollback = tonumber(ARGV[1])
        local stock = tonumber(redis.call("GET", key))
        if stock == nil then
            return -1
        end
        return redis.call("INCRBY", key, rollback)
    `
	result, err := s.rdb.Eval(ctx, luaScript, []string{redisKey}, req.Count).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 执行失败: %v", err)
	}
	// 当返回 -1 时表示 Redis 中库存 key 未初始化，回滚失败
	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "库存回滚失败"}, nil
	}

	// 更新 MySQL 库存，加回扣减的库存
	tx := s.db.Begin()
	if err := tx.Exec("UPDATE product_stocks SET stock = stock + ? WHERE product_id = ?", req.Count, productId).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新 MySQL 库存失败: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	remaining := int32(result.(int64))
	return &stock.StockResp{Code: 0, Message: "库存回滚成功", RemainingStock: &remaining}, nil
}

// ReserveStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ReserveStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// 将产品ID转换为整数
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("产品ID格式错误: %v", err)
	}

	// 统一使用 key "stock:%d" 存储主库存
	redisKey := fmt.Sprintf("stock:%d", productId)
	// 构造用于分布式锁的 key
	lockKey := fmt.Sprintf("lock:%s", redisKey)
	// 获取分布式锁
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second),
		redsync.WithTries(3),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer mutex.UnlockContext(ctx)

	luaScript := `
    local stockKey = KEYS[1]
    local reservedKey = KEYS[2]
    local stock = tonumber(redis.call("GET", stockKey))
    if stock == nil or stock <= 0 then
        return -1
    end
    redis.call("DECRBY", stockKey, 1)
    redis.call("INCRBY", reservedKey, 1)
    return stock - 1
`
	result, err := s.rdb.Eval(ctx, luaScript, []string{redisKey, fmt.Sprintf("reserved:%d", productId)}).Result()
	if err != nil {
		return nil, fmt.Errorf("执行 Redis Lua 脚本失败: %v", err)
	}
	// 当返回 -1 时，说明库存不足或 key 未初始化
	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "库存不足，预占失败"}, nil
	}

	remaining := int32(result.(int64))
	return &stock.StockResp{Code: 0, Message: "库存预占成功", RemainingStock: &remaining}, nil
}

// ReleaseStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ReleaseStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// 将产品 ID 从字符串转换为整数
	productId, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("产品ID格式错误: %v", err)
	}

	// 定义主库存 key 和预留库存 key（注意：两者需要统一管理）
	stockKey := fmt.Sprintf("stock:%d", productId)
	reservedKey := fmt.Sprintf("reserved:%d", productId)
	// 以主库存 key 作为锁定对象，确保库存释放操作的原子性
	lockKey := fmt.Sprintf("lock:%s", stockKey)

	// 获取分布式锁
	mutex := s.redsync.NewMutex(lockKey,
		redsync.WithExpiry(5*time.Second),
		redsync.WithTries(3),
	)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	defer mutex.UnlockContext(ctx)

	luaScript := `
        local stockKey = KEYS[1]
        local reservedKey = KEYS[2]
        local reservedStock = tonumber(redis.call("GET", reservedKey))
        if reservedStock == nil then
            return -1
        end
        if reservedStock <= 0 then
            return 0
        end
        redis.call("DECRBY", reservedKey, 1)
        redis.call("INCRBY", stockKey, 1)
        return 0
    `
	// 执行 Lua 脚本，传入主库存和预留库存的 key
	result, err := s.rdb.Eval(ctx, luaScript, []string{stockKey, reservedKey}).Result()
	if err != nil {
		return nil, fmt.Errorf("执行 Redis Lua 脚本失败: %v", err)
	}
	// 当返回 -1 时，表示 Redis 中库存 key 未初始化，视为错误
	if result.(int64) == -1 {
		return &stock.StockResp{Code: 1, Message: "库存回滚失败"}, nil
	}

	// 同步更新 MySQL 数据库：释放库存意味着主库存加回 1 个单位
	tx := s.db.Begin()
	if err := tx.Exec("UPDATE product_stocks SET stock = stock + ? WHERE product_id = ?", 1, productId).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新 MySQL 库存失败: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	return &stock.StockResp{Code: 0, Message: "库存释放成功"}, nil
}
