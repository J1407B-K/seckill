namespace go stock

include "common.thrift"

// 库存操作请求
struct StockReq {
  1: string productId
  2: i32 count
}

// 库存操作响应
struct StockResp {
  1: i32 code,
  2: string message,
  3: optional i32 remainingStock    // 操作后返回剩余库存
}

service StockService {
  // 查询库存接口
  StockResp QueryStock(1: StockReq req),

  // 预扣库存接口（扣减库存）
  StockResp PreDeductStock(1: StockReq req),

  // 库存回滚接口（库存补偿）
  StockResp RollbackStock(1: StockReq req),

  // 预占库存接口（下单锁库存，但未支付）
  StockResp ReserveStock(1: StockReq req),

  // 释放库存接口（订单取消或支付失败）
  StockResp ReleaseStock(1: StockReq req)
}

